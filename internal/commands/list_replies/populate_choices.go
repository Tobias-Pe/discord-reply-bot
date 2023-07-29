package list_replies

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/cache"
	"github.com/Tobias-Pe/discord-reply-bot/internal/commands"
	"github.com/Tobias-Pe/discord-reply-bot/internal/logger"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func populateChoices(s *discordgo.Session, i *discordgo.InteractionCreate) {

	data := i.ApplicationCommandData()

	var focusedDataOption *discordgo.ApplicationCommandInteractionDataOption
	var nameToOptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	for _, option := range data.Options {
		if option.Focused {
			focusedDataOption = option
		}
		nameToOptionMap[option.Name] = option
	}

	var choices []*discordgo.ApplicationCommandOptionChoice

	switch focusedDataOption.Name {
	case matchTypeOptionName:
		choices = populateChoicesForMatchType(focusedDataOption, choices)
	case toBeMatchedOptionName:
		choices = populateChoicesForToBeMatched(focusedDataOption, choices, nameToOptionMap)
	default:
		logger.Logger.Errorw("Unsupported name for option found", "name", focusedDataOption.Name)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices, // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
		},
	})

	if err != nil {
		logger.Logger.Error("Error on responding with choices", err)
	}
}
func populateChoicesForToBeMatched(focusedDataOption *discordgo.ApplicationCommandInteractionDataOption, choices []*discordgo.ApplicationCommandOptionChoice, nameToOptionMap map[string]*discordgo.ApplicationCommandInteractionDataOption) []*discordgo.ApplicationCommandOptionChoice {
	if focusedDataOption.StringValue() != "" {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  focusedDataOption.StringValue(),
			Value: strings.ToLower(focusedDataOption.StringValue()),
		})
	}

	messageMatches, err := cache.GetAllKeys()
	if err != nil {
		logger.Logger.Debugw("Couldn't fetch all messageMatches", "Error", err)
		return choices
	}

	var keysStringVersion []string
	matchTypeOption := nameToOptionMap[matchTypeOptionName]
	if matchTypeOption != nil {
		for _, messageMatch := range messageMatches {
			isExactMatch := matchTypeOption.StringValue() == models.AllMatchChoices[0]
			if messageMatch.IsExactMatch == isExactMatch {
				keysStringVersion = append(keysStringVersion, messageMatch.Message)
			}
		}
	} else {
		var stringMap = map[string]bool{}
		for _, messageMatch := range messageMatches {
			stringMap[messageMatch.Message] = true
		}
		for key := range stringMap {
			keysStringVersion = append(keysStringVersion, key)
		}
	}

	selectedMatchChoices := commands.SearchChoices(focusedDataOption.StringValue(), keysStringVersion)
	choices = append(choices, commands.TransformSelectedChoices(selectedMatchChoices)...)
	return choices
}

func populateChoicesForMatchType(focusedDataOption *discordgo.ApplicationCommandInteractionDataOption, choices []*discordgo.ApplicationCommandOptionChoice) []*discordgo.ApplicationCommandOptionChoice {
	selectedMatchChoices := commands.SearchChoices(focusedDataOption.StringValue(), models.AllMatchChoices)
	choices = commands.TransformSelectedChoices(selectedMatchChoices)
	return choices
}
