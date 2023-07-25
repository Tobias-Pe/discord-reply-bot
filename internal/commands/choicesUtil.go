package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func TransformSelectedChoices(selectedMatchChoices []string) []*discordgo.ApplicationCommandOptionChoice {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, selectedMatchChoice := range selectedMatchChoices {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  selectedMatchChoice,
			Value: selectedMatchChoice,
		})
	}
	return choices
}

func SearchChoices(searchInput string, allChoices []string) []string {
	var selectedChoices []string
	if len(searchInput) == 0 {
		selectedChoices = allChoices
	} else {
		for _, matchChoice := range allChoices {
			if strings.Contains(matchChoice, searchInput) {
				selectedChoices = append(selectedChoices, matchChoice)
			}
		}
	}
	return selectedChoices
}
