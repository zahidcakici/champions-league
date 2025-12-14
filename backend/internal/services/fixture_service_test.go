package services

import (
	"testing"

	"github.com/zahidcakici/champions-league/internal/models"
)

func TestGenerateSingleRoundRobin(t *testing.T) {
	service := &fixtureService{}

	teams := []models.Team{
		{ID: 1, Name: "Team A", Power: 80},
		{ID: 2, Name: "Team B", Power: 75},
		{ID: 3, Name: "Team C", Power: 85},
		{ID: 4, Name: "Team D", Power: 70},
	}

	matches := service.generateSingleRoundRobin(teams)

	// For 4 teams, single round-robin should produce 6 matches (n*(n-1)/2)
	expectedMatches := 6
	if len(matches) != expectedMatches {
		t.Errorf("Expected %d matches, got %d", expectedMatches, len(matches))
	}

	// Verify all matches have valid team IDs
	for _, match := range matches {
		if match.HomeTeamID == match.AwayTeamID {
			t.Errorf("Match has same home and away team: %d", match.HomeTeamID)
		}
		if match.HomeTeamID < 1 || match.HomeTeamID > 4 {
			t.Errorf("Invalid home team ID: %d", match.HomeTeamID)
		}
		if match.AwayTeamID < 1 || match.AwayTeamID > 4 {
			t.Errorf("Invalid away team ID: %d", match.AwayTeamID)
		}
	}

	// Verify each team plays exactly 3 matches (against each other team once)
	teamMatchCount := make(map[uint]int)
	for _, match := range matches {
		teamMatchCount[match.HomeTeamID]++
		teamMatchCount[match.AwayTeamID]++
	}

	for teamID, count := range teamMatchCount {
		if count != 3 {
			t.Errorf("Team %d should play 3 matches, played %d", teamID, count)
		}
	}
}

func TestGenerateRoundRobin(t *testing.T) {
	service := &fixtureService{}

	teams := []models.Team{
		{ID: 1, Name: "Team A", Power: 80},
		{ID: 2, Name: "Team B", Power: 75},
		{ID: 3, Name: "Team C", Power: 85},
		{ID: 4, Name: "Team D", Power: 70},
	}

	matches := service.generateRoundRobin(teams)

	// For 4 teams, double round-robin should produce 12 matches (n*(n-1))
	expectedMatches := 12
	if len(matches) != expectedMatches {
		t.Errorf("Expected %d matches, got %d", expectedMatches, len(matches))
	}

	// Verify each team plays exactly 6 matches (home and away against each team)
	teamMatchCount := make(map[uint]int)
	for _, match := range matches {
		teamMatchCount[match.HomeTeamID]++
		teamMatchCount[match.AwayTeamID]++
	}

	for teamID, count := range teamMatchCount {
		if count != 6 {
			t.Errorf("Team %d should play 6 matches, played %d", teamID, count)
		}
	}

	// Verify each pair of teams plays exactly 2 matches (home and away)
	pairCount := make(map[string]int)
	for _, match := range matches {
		// Create a unique key for the pair (order doesn't matter)
		var key string
		if match.HomeTeamID < match.AwayTeamID {
			key = string(rune(match.HomeTeamID)) + "-" + string(rune(match.AwayTeamID))
		} else {
			key = string(rune(match.AwayTeamID)) + "-" + string(rune(match.HomeTeamID))
		}
		pairCount[key]++
	}

	for pair, count := range pairCount {
		if count != 2 {
			t.Errorf("Pair %s should play 2 matches, played %d", pair, count)
		}
	}
}

func TestMatchesSpreadAcrossWeeks(t *testing.T) {
	service := &fixtureService{}

	teams := []models.Team{
		{ID: 1, Name: "Team A", Power: 80},
		{ID: 2, Name: "Team B", Power: 75},
		{ID: 3, Name: "Team C", Power: 85},
		{ID: 4, Name: "Team D", Power: 70},
	}

	matches := service.generateRoundRobin(teams)

	// Count matches per week
	weekCount := make(map[int]int)
	for _, match := range matches {
		weekCount[match.Week]++
	}

	// For 4 teams, we should have 6 weeks with 2 matches each
	expectedWeeks := 6
	if len(weekCount) != expectedWeeks {
		t.Errorf("Expected %d weeks, got %d", expectedWeeks, len(weekCount))
	}

	// Each week should have exactly 2 matches (4 teams / 2)
	for week, count := range weekCount {
		if count != 2 {
			t.Errorf("Week %d should have 2 matches, has %d", week, count)
		}
	}
}
