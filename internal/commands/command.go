package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Cmd      *discordgo.ApplicationCommand
	Callback func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
