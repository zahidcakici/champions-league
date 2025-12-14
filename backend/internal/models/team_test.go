package models

import (
	"testing"
)

func TestDefaultTeams(t *testing.T) {
	teams := DefaultTeams()

	if len(teams) != 4 {
		t.Errorf("Expected 4 default teams, got %d", len(teams))
	}

	for _, team := range teams {
		if team.Power < 1 || team.Power > 100 {
			t.Errorf("Team %s has invalid power rating: %d", team.Name, team.Power)
		}
	}

	for _, team := range teams {
		if team.Name == "" {
			t.Error("Team has empty name")
		}
	}

	expectedTeams := map[string]bool{
		"Chelsea":         false,
		"Arsenal":         false,
		"Manchester City": false,
		"Liverpool":       false,
	}

	for _, team := range teams {
		if _, exists := expectedTeams[team.Name]; exists {
			expectedTeams[team.Name] = true
		}
	}

	for name, found := range expectedTeams {
		if !found {
			t.Errorf("Expected team %s not found in default teams", name)
		}
	}
}

func TestTeamPowerRatings(t *testing.T) {
	teams := DefaultTeams()

	var maxPower int
	var maxTeam string
	for _, team := range teams {
		if team.Power > maxPower {
			maxPower = team.Power
			maxTeam = team.Name
		}
	}

	if maxTeam != "Manchester City" {
		t.Logf("Note: Manchester City expected to have highest power, but %s has %d", maxTeam, maxPower)
	}
}
