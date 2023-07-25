package commands

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var allMatchChoices = []string{"exact", "occurrence"}

var addReplyCommand = &discordgo.ApplicationCommand{
	Name:        "add-reply",
	Description: "addReplyCommand for demonstrating options",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "match-type",
			Description:  "How the message should be matched",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "to-be-matched",
			Description: "What needs to be matched",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "to-be-answered",
			Description: "What needs to be answered",
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
		var selectedMatchChoises []string
		userSearchInput := data.Options[0].StringValue()
		if len(userSearchInput) == 0 {
			selectedMatchChoises = allMatchChoices
		} else {
			for _, matchChoice := range allMatchChoices {
				if strings.Contains(matchChoice, userSearchInput) {
					selectedMatchChoises = append(selectedMatchChoises, matchChoice)
				}
			}
		}

		for _, selectedMatchChoise := range selectedMatchChoises {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  cases.Title(language.English).String(selectedMatchChoise),
				Value: selectedMatchChoise,
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

var AddReply = internal.Command{
	Cmd:      addReplyCommand,
	Callback: addReplyFunction,
}
