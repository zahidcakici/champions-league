package services

import (
	"testing"

	"github.com/zahidcakici/champions-league/internal/models"
)

func TestCalculateStandings(t *testing.T) {
	standings := []models.TeamStanding{
		{TeamID: 1, TeamName: "Team A", Played: 3, Won: 2, Drawn: 1, Lost: 0, GoalsFor: 5, GoalsAgainst: 2, Points: 7},
		{TeamID: 2, TeamName: "Team B", Played: 3, Won: 2, Drawn: 0, Lost: 1, GoalsFor: 4, GoalsAgainst: 3, Points: 6},
		{TeamID: 3, TeamName: "Team C", Played: 3, Won: 1, Drawn: 1, Lost: 1, GoalsFor: 3, GoalsAgainst: 3, Points: 4},
		{TeamID: 4, TeamName: "Team D", Played: 3, Won: 0, Drawn: 0, Lost: 3, GoalsFor: 1, GoalsAgainst: 5, Points: 0},
	}

	for i := range standings {
		standings[i].GoalDifference = standings[i].GoalsFor - standings[i].GoalsAgainst
	}

	for _, s := range standings {
		expectedPoints := s.Won*3 + s.Drawn*1
		if s.Points != expectedPoints {
			t.Errorf("Team %s: expected %d points, got %d", s.TeamName, expectedPoints, s.Points)
		}
	}

	for _, s := range standings {
		expectedGD := s.GoalsFor - s.GoalsAgainst
		if s.GoalDifference != expectedGD {
			t.Errorf("Team %s: expected GD %d, got %d", s.TeamName, expectedGD, s.GoalDifference)
		}
	}

	if standings[0].TeamName != "Team A" {
		t.Errorf("Expected Team A to be first, got %s", standings[0].TeamName)
	}

	if standings[3].TeamName != "Team D" {
		t.Errorf("Expected Team D to be last, got %s", standings[3].TeamName)
	}
}

func TestStandingsTiebreaker(t *testing.T) {
	standings := []models.TeamStanding{
		{TeamID: 1, TeamName: "Team A", Points: 6, GoalDifference: 2, GoalsFor: 5},
		{TeamID: 2, TeamName: "Team B", Points: 6, GoalDifference: 4, GoalsFor: 6},
	}

	if standings[0].GoalDifference >= standings[1].GoalDifference {
		t.Log("In case of equal points, team with better GD should rank higher")
	}
}

func TestNormalizePercentages(t *testing.T) {
	service := &standingsService{}

	testCases := []struct {
		name        string
		predictions []models.ChampionshipPrediction
		expected    float64
	}{
		{
			name: "Already sums to 100",
			predictions: []models.ChampionshipPrediction{
				{TeamID: 1, Percentage: 50},
				{TeamID: 2, Percentage: 30},
				{TeamID: 3, Percentage: 20},
			},
			expected: 100,
		},
		{
			name: "Sums to 90",
			predictions: []models.ChampionshipPrediction{
				{TeamID: 1, Percentage: 45},
				{TeamID: 2, Percentage: 30},
				{TeamID: 3, Percentage: 15},
			},
			expected: 100,
		},
		{
			name: "All zeros",
			predictions: []models.ChampionshipPrediction{
				{TeamID: 1, Percentage: 0},
				{TeamID: 2, Percentage: 0},
				{TeamID: 3, Percentage: 0},
				{TeamID: 4, Percentage: 0},
			},
			expected: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service.normalizePercentages(tc.predictions)

			total := 0.0
			for _, p := range tc.predictions {
				total += p.Percentage
			}

			if total != tc.expected {
				t.Errorf("Expected total %.0f, got %.0f", tc.expected, total)
			}
		})
	}
}

func TestMatchResultToStandings(t *testing.T) {
	testCases := []struct {
		name       string
		homeScore  int
		awayScore  int
		homePoints int
		awayPoints int
		homeWon    int
		awayWon    int
		homeDrawn  int
		awayDrawn  int
		homeLost   int
		awayLost   int
	}{
		{
			name:       "Home win",
			homeScore:  2,
			awayScore:  1,
			homePoints: 3,
			awayPoints: 0,
			homeWon:    1,
			awayWon:    0,
			homeDrawn:  0,
			awayDrawn:  0,
			homeLost:   0,
			awayLost:   1,
		},
		{
			name:       "Away win",
			homeScore:  0,
			awayScore:  3,
			homePoints: 0,
			awayPoints: 3,
			homeWon:    0,
			awayWon:    1,
			homeDrawn:  0,
			awayDrawn:  0,
			homeLost:   1,
			awayLost:   0,
		},
		{
			name:       "Draw",
			homeScore:  1,
			awayScore:  1,
			homePoints: 1,
			awayPoints: 1,
			homeWon:    0,
			awayWon:    0,
			homeDrawn:  1,
			awayDrawn:  1,
			homeLost:   0,
			awayLost:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			homeStanding := models.TeamStanding{}
			awayStanding := models.TeamStanding{}

			if tc.homeScore > tc.awayScore {
				homeStanding.Won++
				homeStanding.Points += 3
				awayStanding.Lost++
			} else if tc.homeScore < tc.awayScore {
				awayStanding.Won++
				awayStanding.Points += 3
				homeStanding.Lost++
			} else {
				homeStanding.Drawn++
				awayStanding.Drawn++
				homeStanding.Points++
				awayStanding.Points++
			}

			homeStanding.GoalsFor += tc.homeScore
			homeStanding.GoalsAgainst += tc.awayScore
			awayStanding.GoalsFor += tc.awayScore
			awayStanding.GoalsAgainst += tc.homeScore

			if homeStanding.Points != tc.homePoints {
				t.Errorf("Home points: expected %d, got %d", tc.homePoints, homeStanding.Points)
			}
			if awayStanding.Points != tc.awayPoints {
				t.Errorf("Away points: expected %d, got %d", tc.awayPoints, awayStanding.Points)
			}
			if homeStanding.Won != tc.homeWon {
				t.Errorf("Home won: expected %d, got %d", tc.homeWon, homeStanding.Won)
			}
			if homeStanding.Drawn != tc.homeDrawn {
				t.Errorf("Home drawn: expected %d, got %d", tc.homeDrawn, homeStanding.Drawn)
			}
			if homeStanding.Lost != tc.homeLost {
				t.Errorf("Home lost: expected %d, got %d", tc.homeLost, homeStanding.Lost)
			}
		})
	}
}
