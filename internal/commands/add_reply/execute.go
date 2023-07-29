package add_reply

import (
	"fmt"
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/Tobias-Pe/discord-reply-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
)

func addReply(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	response := fmt.Sprintf(
		"An `%s`-match was added for `%s` with the response `%s`",
		data.Options[0].StringValue(),
		data.Options[1].StringValue(),
		data.Options[2].StringValue(),
	)

	isExactMatch := data.Options[0].StringValue() == models.AllMatchChoices[0]
	err := storage.AddElement(
		models.MessageMatch{
			Message:      data.Options[1].StringValue(),
			IsExactMatch: isExactMatch,
		},
		data.Options[2].StringValue(),
	)
	if err != nil {
		logger.Logger.Error(err)
		response = fmt.Sprintf(
			"Something went wrong. Please check your input.",
		)
	} else {
		logger.Logger.Debugw("Reply Added!", "Match-type", data.Options[0].StringValue(),
			"to-be-responded", data.Options[1].StringValue(),
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
