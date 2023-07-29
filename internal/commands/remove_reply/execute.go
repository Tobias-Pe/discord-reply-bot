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

	response := fmt.Sprintf(
		"The value `%s` for key `%s` got removed",
		data.Options[2].StringValue(),
		data.Options[1].StringValue(),
	)

	toBeRemovedMatch := models.MessageMatch{Message: data.Options[1].StringValue(), IsExactMatch: data.Options[0].StringValue() == models.AllMatchChoices[0]}
	err := storage.RemoveElement(toBeRemovedMatch, data.Options[2].StringValue())
	if err != nil {
		logger.Logger.Warnw("Couldn't delete key-value", "Key", toBeRemovedMatch, "value", data.Options[2].StringValue())
		response = fmt.Sprintf(
			"Something went wrong. Please check your input.",
		)
	} else {
		logger.Logger.Debugw("Reply removed!", "to-be-matched", toBeRemovedMatch,
			"to-be-answered", data.Options[2].StringValue())
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
