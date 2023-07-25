package main

import (
	"flag"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/Tobias-Pe/discord-reply-bot/internal/handler"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers applicationCommands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() {
	flag.Parse()

	if *BotToken == "" {
		value, ok := os.LookupEnv("DISCORD_BOT_TOKEN")
		if !ok {
			panic("No Token set!")
		}
		BotToken = &value
	}

	if *GuildID == "" {
		value, ok := os.LookupEnv("DISCORD_GUILD_ID")
		if !ok {
			panic("No GuildID set!")
		}
		GuildID = &value
	}
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	applicationCommands = []*discordgo.ApplicationCommand{
		commands.AddReply.Cmd,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		commands.AddReply.Cmd.Name: commands.AddReply.Callback,
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(handler.MessageCreate)

	openSession()

	addCommands()

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			return
		}
	}(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}

func openSession() {
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

func addCommands() {

	log.Println("Deleting applicationCommands...")
	cmds, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	if err == nil {
		for _, cmd := range cmds {
			_ = s.ApplicationCommandDelete(s.State.User.ID, *GuildID, cmd.ID)
		}
	}

	log.Println("Adding applicationCommands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(applicationCommands))
	for i, v := range applicationCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}
