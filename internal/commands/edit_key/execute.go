package edit_key

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
)

func editKey(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var nameToOptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range data.Options {
		nameToOptionMap[option.Name] = option
	}

	response := fmt.Sprintf(
		"An `%s`-match got changed from `%s` to `%s`",
		nameToOptionMap[matchTypeOptionName].StringValue(),
		nameToOptionMap[toBeMatchedOptionName].StringValue(),
		nameToOptionMap[newToBeMatchedOptionName].StringValue(),
	)

	isExactMatch := nameToOptionMap[matchTypeOptionName].StringValue() == models.AllMatchChoices[0]
	err := storage.RenameKey(
		models.MessageMatch{
			Message:      nameToOptionMap[toBeMatchedOptionName].StringValue(),
			IsExactMatch: isExactMatch,
		},
		models.MessageMatch{
			Message:      nameToOptionMap[newToBeMatchedOptionName].StringValue(),
			IsExactMatch: isExactMatch,
		},
	)
	if err != nil {
		logger.Logger.Error(err)
		response = fmt.Sprintf(
			"Something went wrong. Please check your input.",
		)
	} else {
		logger.Logger.Debugw("Key replaced!", "Match-type", nameToOptionMap[matchTypeOptionName].StringValue(),
			"to-be-responded", nameToOptionMap[toBeMatchedOptionName].StringValue(),
			"new-to-be-answered", nameToOptionMap[newToBeMatchedOptionName].StringValue())

		cache.InvalidateKeyCache()
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: response,
		},
	})

	if err != nil {
		logger.Logger.Error("Error on responding", err)
	}
}
