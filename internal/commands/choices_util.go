package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func TransformSelectedChoices(selectedMatchChoices []string) []*discordgo.ApplicationCommandOptionChoice {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, selectedMatchChoice := range selectedMatchChoices {
		if len(selectedMatchChoice) <= 100 {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  selectedMatchChoice,
				Value: selectedMatchChoice,
			})
		}
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
func SearchChoicesIndices(searchInput string, allChoices []string) []int {
	var selectedChoices []int
	for i, matchChoice := range allChoices {
		if len(searchInput) == 0 || strings.Contains(matchChoice, searchInput) {
			selectedChoices = append(selectedChoices, i)
		}
	}
	return selectedChoices
}
