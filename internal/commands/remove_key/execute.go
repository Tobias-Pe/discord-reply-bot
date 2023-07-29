package remove_key

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
)

func removeKey(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var nameToOptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range data.Options {
		nameToOptionMap[option.Name] = option
	}

	response := fmt.Sprintf(
		"The key `%s` got removed",
		nameToOptionMap[toBeMatchedOptionName].StringValue(),
	)

	toBeRemovedMatch := models.MessageMatch{Message: nameToOptionMap[toBeMatchedOptionName].StringValue(), IsExactMatch: nameToOptionMap[matchTypeOptionName].StringValue() == models.AllMatchChoices[0]}
	err := storage.RemoveKey(toBeRemovedMatch)
	if err != nil {
		logger.Logger.Warnw("Couldn't delete key", "Key", toBeRemovedMatch)
		response = fmt.Sprintf(
			"Something went wrong. Please check your input.",
		)
	} else {
		logger.Logger.Debugw("Key removed!", "to-be-matched", toBeRemovedMatch)
		cache.InvalidateKeyCache()
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: response,
		},
	})

	if err != nil {
		logger.Logger.Error("Error on responding", err)
	}
}
