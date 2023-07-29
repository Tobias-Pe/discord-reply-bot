package remove_key

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

const matchTypeOptionName = "match-type"
const toBeMatchedOptionName = "to-be-matched"

var removeKeyCommand = &discordgo.ApplicationCommand{
	Name:        "remove-key",
	Description: "Remove a key",
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
			Description:  "Select from which match you want to delete.",
			Required:     true,
			Autocomplete: true,
		},
	},
}

var removeKeyFunction = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		removeKey(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		populateChoices(s, i)
	}
}

var RemoveKey = commands.Command{
	Cmd:      removeKeyCommand,
	Callback: removeKeyFunction,
}
