package edit_reply

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

const matchTypeOptionName = "match-type"
const toBeMatchedOptionName = "to-be-matched"
const toBeAnsweredOptionName = "to-be-answered"
const newToBeAnsweredOptionName = "new-to-be-answered"

var editReplyCommand = &discordgo.ApplicationCommand{
	Name:        "edit-reply",
	Description: "Edit a reply to a strict or not so strict matched message",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         matchTypeOptionName,
			Description:  "Which type of message you want to edit.",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         toBeMatchedOptionName,
			Description:  "Which match you want to edit.",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         toBeAnsweredOptionName,
			Description:  "Which answer you want to edit.",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        newToBeAnsweredOptionName,
			Description: "What do you want as an answer instead?",
			Required:    true,
		},
	},
}

var editReplyFunction = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		editReply(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		populateChoices(s, i)
	}
}

var EditReply = commands.Command{
	Cmd:      editReplyCommand,
	Callback: editReplyFunction,
}
