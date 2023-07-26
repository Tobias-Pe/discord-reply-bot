package remove_reply

import (
	"encoding/json"
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var removeReplyCommand = &discordgo.ApplicationCommand{
	Name:        "remove-reply",
	Description: "Remove a reply",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "to-be-matched",
			Description:  "Select from which match you want to delete an entry.",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "to-be-answered",
			Description:  "Select which answer you want to delete.",
			Required:     true,
			Autocomplete: true,
		},
	},
}

var removeReplyFunction = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		removeReply(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		populateChoices(s, i)
	}
}

func removeReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	var currentSelectedKey models.MessageMatch
	err := json.Unmarshal([]byte(data.Options[0].StringValue()), &currentSelectedKey)
	if err != nil {
		logger.Logger.Debugw("Couldn't unmarshall selected key", "Key", data.Options[0].StringValue())
		return
	}
	err = storage.RemoveElement(currentSelectedKey, data.Options[1].StringValue())
	if err != nil {
		logger.Logger.Warnw("Couldn't delete key-value", "Key", currentSelectedKey, "value", data.Options[1].StringValue())
		return
	}

	logger.Logger.Debugw("Reply removed!", "to-be-matched", data.Options[0].StringValue(),
		"to-be-answered", data.Options[1].StringValue())

	cache.InvalidateKeyCache()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(
				"The value `%s` for key `%s` got removed",
				data.Options[1].StringValue(),
				currentSelectedKey.Message,
			),
		},
	})

	if err != nil {
		logger.Logger.Error("Error on responding", err)
	}
}

func populateChoices(s *discordgo.Session, i *discordgo.InteractionCreate) {

	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

	switch {
	case data.Options[0].Focused:
		keys, err := cache.GetAllKeys()
		if err != nil {
			logger.Logger.Debugw("Couldn't fetch all keys", "Error", err)
			break
		}

		var keysStringVersion []string
		for _, key := range keys {
			sprintf := fmt.Sprintf("Exact: %t, Msg: %s", key.IsExactMatch, key.Message)
			logger.Logger.Debugf("Key stingified: %s", sprintf)
			keysStringVersion = append(keysStringVersion, sprintf)
		}

		selectedMatchChoiceIndices := commands.SearchChoicesIndices(data.Options[0].StringValue(), keysStringVersion)
		logger.Logger.Debug("Selected Choices: ", selectedMatchChoiceIndices)

		for _, index := range selectedMatchChoiceIndices {
			marshal, err := json.Marshal(keys[index])
			if err != nil {
				logger.Logger.Debugw("Couldn't marshall key", "Key", keys[index])
				choices = []*discordgo.ApplicationCommandOptionChoice{}
				break
			}

			if len(keysStringVersion[index]) <= 100 && len(string(marshal)) <= 100 {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  keysStringVersion[index],
					Value: string(marshal),
				})
			}
		}
		if data.Options[0].StringValue() != "" {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Options[0].StringValue(),
				Value: strings.ToLower(data.Options[0].StringValue()),
			})
		}

		//TODO add case custom key chosen and request to remove
	case data.Options[1].Focused:
		var currentSelectedKey models.MessageMatch
		err := json.Unmarshal([]byte(data.Options[0].StringValue()), &currentSelectedKey)
		if err != nil {
			logger.Logger.Debugw("Couldn't unmarshall selected key", "Key", data.Options[0].StringValue())
			break
		}
		values, err := storage.GetAll(currentSelectedKey)
		if err != nil {
			logger.Logger.Debugw("Couldn't get from redis selected key's values", "Key", currentSelectedKey)
			break
		}
		selectedMatchChoices := commands.SearchChoices(data.Options[1].StringValue(), values)
		choices = commands.TransformSelectedChoices(selectedMatchChoices)
		if data.Options[1].StringValue() != "" {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Options[1].StringValue(),
				Value: strings.ToLower(data.Options[1].StringValue()),
			})
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices, // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
		},
	})

	if err != nil {
		logger.Logger.Error("Error on responding with choices", err)
	}
}

var RemoveReply = commands.Command{
	Cmd:      removeReplyCommand,
	Callback: removeReplyFunction,
}
