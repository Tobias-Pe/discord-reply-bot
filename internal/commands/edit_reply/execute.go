package edit_reply

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
)

func editReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var nameToOptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range data.Options {
		nameToOptionMap[option.Name] = option
	}

	response := fmt.Sprintf(
		"An `%s`-match was edited for `%s` with the response `%s` got changed to `%s`",
		nameToOptionMap[matchTypeOptionName].StringValue(),
		nameToOptionMap[toBeMatchedOptionName].StringValue(),
		nameToOptionMap[toBeAnsweredOptionName].StringValue(),
		nameToOptionMap[newToBeAnsweredOptionName].StringValue(),
	)

	isExactMatch := nameToOptionMap[matchTypeOptionName].StringValue() == models.AllMatchChoices[0]
	toBeRemovedMatch := models.MessageMatch{Message: nameToOptionMap[toBeMatchedOptionName].StringValue(), IsExactMatch: isExactMatch}
	err := storage.RemoveElement(toBeRemovedMatch, nameToOptionMap[toBeAnsweredOptionName].StringValue())
	if err != nil {
		logger.Logger.Warnw("Couldn't delete key-value", "Key", toBeRemovedMatch, "value", nameToOptionMap[toBeAnsweredOptionName].StringValue())
		response = fmt.Sprintf(
			"Something went wrong. Please check your input.",
		)
	} else {
		err = storage.AddElement(
			models.MessageMatch{
				Message:      nameToOptionMap[toBeMatchedOptionName].StringValue(),
				IsExactMatch: isExactMatch,
			},
			nameToOptionMap[newToBeAnsweredOptionName].StringValue(),
		)
		if err != nil {
			logger.Logger.Error(err)
			response = fmt.Sprintf(
				"Something went wrong. Please check your input.",
			)
		} else {
			logger.Logger.Debugw("Reply replaced!", "Match-type", nameToOptionMap[matchTypeOptionName].StringValue(),
				"to-be-responded", nameToOptionMap[toBeMatchedOptionName].StringValue(),
				"to-be-answered", nameToOptionMap[toBeAnsweredOptionName].StringValue(),
				"new-to-be-answered", nameToOptionMap[newToBeAnsweredOptionName].StringValue())

			cache.InvalidateKeyCache()
		}
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
