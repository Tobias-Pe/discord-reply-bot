package remove_reply

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
)

func removeReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var nameToOptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range data.Options {
		nameToOptionMap[option.Name] = option
	}

	response := fmt.Sprintf(
		"The value `%s` for key `%s` got removed",
		nameToOptionMap[toBeAnsweredOptionName].StringValue(),
		nameToOptionMap[toBeMatchedOptionName].StringValue(),
	)

	toBeRemovedMatch := models.MessageMatch{Message: nameToOptionMap[toBeMatchedOptionName].StringValue(), IsExactMatch: nameToOptionMap[matchTypeOptionName].StringValue() == models.AllMatchChoices[0]}
	err := storage.RemoveElement(toBeRemovedMatch, nameToOptionMap[toBeAnsweredOptionName].StringValue())
	if err != nil {
		logger.Logger.Warnw("Couldn't delete key-value", "Key", toBeRemovedMatch, "value", nameToOptionMap[toBeAnsweredOptionName].StringValue())
		response = fmt.Sprintf(
			"Something went wrong. Please check your input.",
		)
	} else {
		logger.Logger.Debugw("Reply removed!", "to-be-matched", toBeRemovedMatch,
			"to-be-answered", nameToOptionMap[toBeAnsweredOptionName].StringValue())
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
