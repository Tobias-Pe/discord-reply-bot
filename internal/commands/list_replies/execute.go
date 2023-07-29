package list_replies

import (
	"encoding/json"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func listReplies(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error
	data := i.ApplicationCommandData()

	var nameToOptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range data.Options {
		nameToOptionMap[option.Name] = option
	}

	var currentSelectionKeys []models.MessageMatch
	if nameToOptionMap[toBeMatchedOptionName] != nil {
		currentSelectionKeys = append(currentSelectionKeys, models.MessageMatch{
			Message:      nameToOptionMap[toBeMatchedOptionName].StringValue(),
			IsExactMatch: false,
		})
		currentSelectionKeys = append(currentSelectionKeys, models.MessageMatch{
			Message:      nameToOptionMap[toBeMatchedOptionName].StringValue(),
			IsExactMatch: true,
		})
	} else {
		currentSelectionKeys, err = cache.GetAllKeys()
		if err != nil {
			logger.Logger.Errorw("Couldnt get all keys", "error", err)
		}
	}

	if nameToOptionMap[matchTypeOptionName] != nil {
		var tmpSelectionKeys []models.MessageMatch
		isExactMatch := nameToOptionMap[matchTypeOptionName].StringValue() == models.AllMatchChoices[0]
		for _, key := range currentSelectionKeys {
			if key.IsExactMatch == isExactMatch {
				tmpSelectionKeys = append(tmpSelectionKeys, key)
			}
		}
		currentSelectionKeys = tmpSelectionKeys
	}

	var selectedCompleteData []models.CompleteData
	for _, key := range currentSelectionKeys {
		var values []string
		values, err = storage.GetAll(key)
		if values != nil && len(values) > 0 {
			selectedCompleteData = append(selectedCompleteData, models.CompleteData{
				MessageMatch: key,
				Reply:        values,
			})
		}
	}

	jsonCompleteData, err := json.MarshalIndent(selectedCompleteData, "", "    ")

	content := "Here you go!"

	if err != nil {
		content = "Something went wrong :("
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Files: []*discordgo.File{
				{
					ContentType: "application/json",
					Name:        "replies_query.json",
					Reader:      strings.NewReader(string(jsonCompleteData)),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Logger.Error("Error on responding", err)
	}
}
