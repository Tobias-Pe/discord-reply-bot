package remove_reply

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

const matchTypeOptionName = "match-type"
const toBeMatchedOptionName = "to-be-matched"
const toBeAnsweredOptionName = "to-be-answered"

var removeReplyCommand = &discordgo.ApplicationCommand{
	Name:        "remove-reply",
	Description: "Remove a reply",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         matchTypeOptionName,
			Description:  "Select which type of match you want to delete.",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         toBeMatchedOptionName,
			Description:  "Select from which match you want to delete an entry.",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         toBeAnsweredOptionName,
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

var RemoveReply = commands.Command{
	Cmd:      removeReplyCommand,
	Callback: removeReplyFunction,
}
