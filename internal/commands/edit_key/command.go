package edit_key

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

const matchTypeOptionName = "match-type"
const toBeMatchedOptionName = "to-be-matched"
const newToBeMatchedOptionName = "new-to-be-matched"

var editKeyCommand = &discordgo.ApplicationCommand{
	Name:        "edit-key",
	Description: "Edit a key to a strict or not so strict matched message",
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
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        newToBeMatchedOptionName,
			Description: "What do you want as an answer instead?",
			Required:    true,
		},
	},
}

var editKeyFunction = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		editKey(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		populateChoices(s, i)
	}
}

var EditKey = commands.Command{
	Cmd:      editKeyCommand,
	Callback: editKeyFunction,
}
