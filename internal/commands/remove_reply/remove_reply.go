package remove_reply

import (
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
			Name:         "match-type",
			Description:  "Select which type of match you want to delete.",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
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

	toBeRemovedMatch := models.MessageMatch{Message: data.Options[1].StringValue(), IsExactMatch: data.Options[0].StringValue() == models.AllMatchChoices[0]}
	err := storage.RemoveElement(toBeRemovedMatch, data.Options[2].StringValue())
	if err != nil {
		logger.Logger.Warnw("Couldn't delete key-value", "Key", toBeRemovedMatch, "value", data.Options[2].StringValue())
		return
	}

	logger.Logger.Debugw("Reply removed!", "to-be-matched", toBeRemovedMatch,
		"to-be-answered", data.Options[2].StringValue())

	cache.InvalidateKeyCache()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(
				"The value `%s` for key `%s` got removed",
				data.Options[2].StringValue(),
				toBeRemovedMatch.Message,
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
		selectedMatchChoices := commands.SearchChoices(data.Options[0].StringValue(), models.AllMatchChoices)
		choices = commands.TransformSelectedChoices(selectedMatchChoices)
	case data.Options[1].Focused:
		if data.Options[1].StringValue() != "" {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Options[1].StringValue(),
				Value: strings.ToLower(data.Options[1].StringValue()),
			})
		}

		messageMatches, err := cache.GetAllKeys()
		if err != nil {
			logger.Logger.Debugw("Couldn't fetch all messageMatches", "Error", err)
			break
		}

		var keysStringVersion []string
		for _, messageMatch := range messageMatches {
			isExactMatch := data.Options[0].StringValue() == models.AllMatchChoices[0]
			if messageMatch.IsExactMatch == isExactMatch {
				keysStringVersion = append(keysStringVersion, messageMatch.Message)
			}
		}

		selectedMatchChoices := commands.SearchChoices(data.Options[1].StringValue(), keysStringVersion)
		choices = append(choices, commands.TransformSelectedChoices(selectedMatchChoices)...)

	case data.Options[2].Focused:
		if data.Options[2].StringValue() != "" {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Options[2].StringValue(),
				Value: strings.ToLower(data.Options[2].StringValue()),
			})
		}

		isExactMatch := data.Options[0].StringValue() == models.AllMatchChoices[0]
		values, err := storage.GetAll(models.MessageMatch{Message: data.Options[1].StringValue(), IsExactMatch: isExactMatch})
		if err != nil {
			logger.Logger.Debugw("Couldn't get key's values", "Key", data.Options[0].StringValue(), "Excact", true)
			break
		}

		selectedMatchChoices := commands.SearchChoices(data.Options[2].StringValue(), values)
		choices = append(choices, commands.TransformSelectedChoices(selectedMatchChoices)...)
	}

	logger.Logger.Debug(choices)
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
