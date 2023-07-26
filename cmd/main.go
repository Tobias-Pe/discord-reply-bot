package main

import (
	"flag"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/add-reply"
	"github.com/Tobias-Pe/discord-reply-bot/internal/handler/messages"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers applicationCommands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() {
	logger.InitLogger()

	flag.Parse()

	if *BotToken == "" {
		value, ok := os.LookupEnv("DISCORD_BOT_TOKEN")
		if !ok {
			logger.Logger.Panic("No Token set!")
		}
		BotToken = &value
	}

	if *GuildID == "" {
		value, ok := os.LookupEnv("DISCORD_GUILD_ID")
		if !ok {
			logger.Logger.Panic("No GuildID set!")
		}
		GuildID = &value
	}
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		logger.Logger.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	applicationCommands = []*discordgo.ApplicationCommand{
		add_reply.AddReply.Cmd,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		add_reply.AddReply.Cmd.Name: add_reply.AddReply.Callback,
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
		logger.Logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Register the messageCreate func as a callback for MessageCreate events.
	s.AddHandler(messages.MessageCreate)

	openSession()

	addCommands()

	defer func(s *discordgo.Session) {
		_ = s.Close()
	}(s)

	defer func(Logger *zap.SugaredLogger) {
		_ = Logger.Sync()
	}(logger.Logger)

	err := storage.Test()
	if err != nil {
		logger.Logger.Panic("Redis not reachable!")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger.Logger.Infof("Press Ctrl+C to exit")
	<-stop

	logger.Logger.Infof("Gracefully shutting down.")
}

func openSession() {
	err := s.Open()
	if err != nil {
		logger.Logger.Fatalf("Cannot open the session: %v", err)
	}
}

func addCommands() {

	logger.Logger.Info("Deleting applicationCommands...")
	cmds, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	if err == nil {
		for _, cmd := range cmds {
			_ = s.ApplicationCommandDelete(s.State.User.ID, *GuildID, cmd.ID)
		}
	}

	logger.Logger.Infof("Adding applicationCommands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(applicationCommands))
	for i, v := range applicationCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			logger.Logger.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}
