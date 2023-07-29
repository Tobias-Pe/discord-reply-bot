package add_reply

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

const matchTypeOptionName = "match-type"
const toBeMatchedOptionName = "to-be-matched"
const toBeAnsweredOptionName = "to-be-answered"

var addReplyCommand = &discordgo.ApplicationCommand{
	Name:        "add-reply",
	Description: "Add a reply to a strict or not so strict matched message",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         matchTypeOptionName,
			Description:  "How the message should be matched. Choose one of the options!",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         toBeMatchedOptionName,
			Description:  "What needs to be matched. Choose one of the last few messages or write your own!",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        toBeAnsweredOptionName,
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

var AddReply = commands.Command{
	Cmd:      addReplyCommand,
	Callback: addReplyFunction,
}
