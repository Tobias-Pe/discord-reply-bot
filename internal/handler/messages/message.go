package messages

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
)

const messagesToKeepInLastMessages = 10

var lastMessages []string

var cacheNeedsUpdate = true
var cachedKeys []models.MessageMatch

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	transformedInput := strings.ToLower(strings.TrimSpace(m.Content))

	updateLastMessages(transformedInput)

	allKeys, err := getAllKeys()

	if err != nil {
		logger.Logger.Warnf("Couldn't fetch all keys %s", err)
		return
	}

	possibleResponses := getPossibleResponses(allKeys, transformedInput)

	respondOnDiscord(s, m, possibleResponses)
}

func respondOnDiscord(s *discordgo.Session, m *discordgo.MessageCreate, possibleResponses []string) {
	if possibleResponses != nil && len(possibleResponses) > 0 {
		randomIndex := rand.Intn(len(possibleResponses))
		_, _ = s.ChannelMessageSend(m.ChannelID, possibleResponses[randomIndex])
	}
}

func getPossibleResponses(allKeys []models.MessageMatch, transformedInput string) []string {
	var possibleResponses []string
	for _, key := range allKeys {
		if key.IsExactMatch && transformedInput == key.Message || !key.IsExactMatch && strings.Contains(transformedInput, key.Message) {
			keysResponses, err := storage.GetAll(key)
			if err != nil {
				logger.Logger.Warnw("Couldn't fetch all values from Redis", "Error", err, "Key", key)
				return nil
			}
			possibleResponses = append(possibleResponses, keysResponses...)
		}
	}
	return possibleResponses
}

func updateLastMessages(transformedInput string) {
	lastMessages = append(lastMessages, transformedInput)
	if len(lastMessages) > messagesToKeepInLastMessages {
		lastMessages = lastMessages[1:]
	}
}

func GetLastMessages() []string {
	return lastMessages
}

func InvalidateKeyCache() {
	cacheNeedsUpdate = true
	logger.Logger.Debug("Cache Invalidated")
}

func getAllKeys() ([]models.MessageMatch, error) {
	var err error
	if cacheNeedsUpdate {
		cachedKeys, err = storage.GetAllKeys()
		if err == nil {
			cacheNeedsUpdate = false
			logger.Logger.Debug("Cache Updated")
		}
	}
	return cachedKeys, err
}
