package messages

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const messagesToKeepInLastMessages = 10

var LastMessages []string

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	transformedInput := strings.ToLower(strings.TrimSpace(m.Content))

	LastMessages = append(LastMessages, transformedInput)
	if len(LastMessages) > messagesToKeepInLastMessages {
		LastMessages = LastMessages[1:]
	}

	answer := storage.GetRandom(transformedInput)

	if len(answer) > 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, answer)
	}
}
