package list_replies

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

const matchTypeOptionName = "match-type"
const toBeMatchedOptionName = "to-be-matched"

var listRepliesCommand = &discordgo.ApplicationCommand{
	Name:        "list-replies",
	Description: "List all or some replies",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         matchTypeOptionName,
			Description:  "Select which type of match you want to list.",
			Autocomplete: true,
		},
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         toBeMatchedOptionName,
			Description:  "Select which to-be-matched word you want to list.",
			Autocomplete: true,
		},
	},
}

var listRepliesFunction = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		listReplies(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		populateChoices(s, i)
	}
}

var ListReplies = commands.Command{
	Cmd:      listRepliesCommand,
	Callback: listRepliesFunction,
}
