package models

import "testing"

func TestZombie(t *testing.T) {
	t.Run("Zombie test allmatches array", func(t *testing.T) {
		if AllMatchChoices[0] != "exact" || len(AllMatchChoices) != 2 {
			t.Errorf("This code need to be adjusted analogue to the business logic! The logic currently relies on it having 2 elements and the first one beeing exact")
		}
	})
}
