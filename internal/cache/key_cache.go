package cache

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
)

var cacheNeedsUpdate = true
var cachedKeys []models.MessageMatch

func GetAllKeys() ([]models.MessageMatch, error) {
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

func InvalidateKeyCache() {
	cacheNeedsUpdate = true
	logger.Logger.Debug("Cache Invalidated")
}
