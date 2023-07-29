package main

import (
	"flag"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/add_reply"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/edit_key"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/edit_reply"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/list_replies"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/remove_key"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands/remove_reply"
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
	RedisUrl = flag.String("redis", "", "Redis url")
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

	if *RedisUrl == "" {
		value, ok := os.LookupEnv("REDIS_URL")
		if !ok {
			logger.Logger.Panic("No RedisUrl set!")
		}
		RedisUrl = &value
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
		remove_key.RemoveKey.Cmd,
		remove_reply.RemoveReply.Cmd,
		list_replies.ListReplies.Cmd,
		edit_reply.EditReply.Cmd,
		edit_key.EditKey.Cmd,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		add_reply.AddReply.Cmd.Name:       add_reply.AddReply.Callback,
		remove_key.RemoveKey.Cmd.Name:     remove_key.RemoveKey.Callback,
		remove_reply.RemoveReply.Cmd.Name: remove_reply.RemoveReply.Callback,
		list_replies.ListReplies.Cmd.Name: list_replies.ListReplies.Callback,
		edit_reply.EditReply.Cmd.Name:     edit_reply.EditReply.Callback,
		edit_key.EditKey.Cmd.Name:         edit_key.EditKey.Callback,
	}
)

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	// Register the command handlers.
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	// Register the message handler.
	s.AddHandler(messages.MessageCreate)

	openSession()

	addCommands()

	defer func(s *discordgo.Session) {
		_ = s.Close()
	}(s)

	defer func(Logger *zap.SugaredLogger) {
		_ = Logger.Sync()
	}(logger.Logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger.Logger.Infof("Press Ctrl+C to exit")
	<-stop

	logger.Logger.Infof("Gracefully shutting down.")
}

func openSession() {
	err := s.Open()
	if err != nil {
		logger.Logger.Panicf("Cannot open the session: %v", err)
	}

	storage.InitClient(*RedisUrl)

	err = storage.CheckConnection()
	if err != nil {
		logger.Logger.Panicf("Redis not reachable at %s: %v", *RedisUrl, err)
	}
}

func addCommands() {
	logger.Logger.Infof("Adding applicationCommands...")
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, *GuildID, applicationCommands)
	if err != nil {
		logger.Logger.Panicf("Couldn't create all commands: %v", err)
	}
}
