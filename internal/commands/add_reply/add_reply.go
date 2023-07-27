package add_reply

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/Tobias-Pe/discord-reply-bot/internal/handler/messages"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var addReplyCommand = &discordgo.ApplicationCommand{
	Name:        "add-reply",
	Description: "Add a reply to a strict or not so strict matched message",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "match-type",
			Description:  "How the message should be matched. Choose one of the options!",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "to-be-matched",
			Description:  "What needs to be matched. Choose one of the last few messages or write your own!",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "to-be-answered",
			Description: "What needs to be answered.",
			Required:    true,
		},
	},
}

var addReplyFunction = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		addReply(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		populateChoices(s, i)
	}
}

func addReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	isExactMatch := data.Options[0].StringValue() == models.AllMatchChoices[0]
	err := storage.AddElement(
		models.MessageMatch{
			Message:      data.Options[1].StringValue(),
			IsExactMatch: isExactMatch,
		},
		data.Options[2].StringValue(),
	)
	if err != nil {
		logger.Logger.Error(err)
		return
	}

	logger.Logger.Debugw("Reply Added!", "Match-type", data.Options[0].StringValue(),
		"to-be-responded", data.Options[1].StringValue(),
		"to-be-answered", data.Options[2].StringValue())

	cache.InvalidateKeyCache()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(
				"An `%s`-match was added for `%s` with the response `%s`",
				data.Options[0].StringValue(),
				data.Options[1].StringValue(),
				data.Options[2].StringValue(),
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
		selectedMatchChoices := commands.SearchChoices(data.Options[1].StringValue(), messages.GetLastMessages())
		choices = append(choices, commands.TransformSelectedChoices(selectedMatchChoices)...)
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

var AddReply = commands.Command{
	Cmd:      addReplyCommand,
	Callback: addReplyFunction,
}
