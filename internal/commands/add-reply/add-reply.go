package add_reply

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/Tobias-Pe/discord-reply-bot/internal/handler/messages"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var allMatchChoices = []string{"exact", "occurrence"}

var addReplyCommand = &discordgo.ApplicationCommand{
	Name:        "add-reply",
	Description: "addReplyCommand for demonstrating options",
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

	storage.AddElement(data.Options[1].StringValue(), data.Options[2].StringValue())

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"An `%s`-match was added for `%s` with the response `%s`",
				data.Options[0].StringValue(),
				data.Options[1].StringValue(),
				data.Options[2].StringValue(),
			),
		},
	})
}

func populateChoices(s *discordgo.Session, i *discordgo.InteractionCreate) {

	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

	switch {
	case data.Options[0].Focused:
		selectedMatchChoices := commands.SearchChoices(data.Options[0].StringValue(), allMatchChoices)
		choices = commands.TransformSelectedChoices(selectedMatchChoices)
	case data.Options[1].Focused:
		selectedMatchChoices := commands.SearchChoices(data.Options[1].StringValue(), messages.LastMessages)
		choices = commands.TransformSelectedChoices(selectedMatchChoices)
		if data.Options[1].StringValue() != "" {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Options[1].StringValue(),
				Value: strings.ToLower(data.Options[1].StringValue()),
			})
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices, // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
		},
	})
}

var AddReply = commands.Command{
	Cmd:      addReplyCommand,
	Callback: addReplyFunction,
}
