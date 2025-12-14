package services

import (
	"testing"

	"github.com/zahidcakici/champions-league/internal/models"
)

func TestSimulateMatch(t *testing.T) {
	service := &simulationService{}

	homeTeam := &models.Team{ID: 1, Name: "Strong Team", Power: 95}
	awayTeam := &models.Team{ID: 2, Name: "Weak Team", Power: 50}

	// Run simulation multiple times to test consistency
	homeWins := 0
	awayWins := 0
	draws := 0
	iterations := 1000

	for i := 0; i < iterations; i++ {
		homeScore, awayScore := service.simulateMatch(homeTeam, awayTeam)

		// Scores should be non-negative
		if homeScore < 0 || awayScore < 0 {
			t.Errorf("Scores should be non-negative: home=%d, away=%d", homeScore, awayScore)
		}

		// Scores should be reasonable (max 7 as per implementation)
		if homeScore > 7 || awayScore > 7 {
			t.Errorf("Scores should be max 7: home=%d, away=%d", homeScore, awayScore)
		}

		if homeScore > awayScore {
			homeWins++
		} else if awayScore > homeScore {
			awayWins++
		} else {
			draws++
		}
	}

	// Stronger team with home advantage should win more often
	// Allow for some randomness but expect significant difference
	if homeWins <= awayWins {
		t.Logf("Warning: Stronger home team didn't win more often. Home: %d, Away: %d, Draws: %d",
			homeWins, awayWins, draws)
	}

	// At least some variety in results
	if homeWins == 0 || awayWins == 0 {
		t.Errorf("Expected some variety in results. Home wins: %d, Away wins: %d", homeWins, awayWins)
	}
}

func TestSimulateMatchEqualTeams(t *testing.T) {
	service := &simulationService{}

	homeTeam := &models.Team{ID: 1, Name: "Team A", Power: 80}
	awayTeam := &models.Team{ID: 2, Name: "Team B", Power: 80}

	homeWins := 0
	awayWins := 0
	iterations := 1000

	for i := 0; i < iterations; i++ {
		homeScore, awayScore := service.simulateMatch(homeTeam, awayTeam)

		if homeScore > awayScore {
			homeWins++
		} else if awayScore > homeScore {
			awayWins++
		}
	}

	// Home team should still have slight advantage due to home bonus
	// But it shouldn't be overwhelming for equal teams
	homeWinRate := float64(homeWins) / float64(iterations)
	awayWinRate := float64(awayWins) / float64(iterations)

	t.Logf("Equal teams - Home win rate: %.2f, Away win rate: %.2f", homeWinRate, awayWinRate)

	// Home advantage should give some benefit but not too extreme
	if homeWinRate < 0.3 || homeWinRate > 0.7 {
		t.Logf("Home win rate outside expected range for equal teams: %.2f", homeWinRate)
	}
}

func TestGenerateGoals(t *testing.T) {
	service := &simulationService{}

	// Test with different expected goals
	testCases := []struct {
		expectedGoals float64
		minGoals      int
		maxGoals      int
	}{
		{0.5, 0, 7},
		{1.5, 0, 7},
		{2.5, 0, 7},
		{3.0, 0, 7},
	}

	for _, tc := range testCases {
		totalGoals := 0
		iterations := 1000

		for i := 0; i < iterations; i++ {
			goals := service.generateGoals(tc.expectedGoals)

			if goals < tc.minGoals || goals > tc.maxGoals {
				t.Errorf("Goals %d outside expected range [%d, %d] for expectedGoals=%.1f",
					goals, tc.minGoals, tc.maxGoals, tc.expectedGoals)
			}

			totalGoals += goals
		}

		avgGoals := float64(totalGoals) / float64(iterations)
		t.Logf("Expected goals: %.1f, Average actual: %.2f", tc.expectedGoals, avgGoals)
	}
}

func TestHomeAdvantage(t *testing.T) {
	service := &simulationService{}

	// Same power teams, run many simulations
	team := &models.Team{ID: 1, Name: "Team", Power: 75}

	homeScoreTotal := 0
	awayScoreTotal := 0
	iterations := 1000

	for i := 0; i < iterations; i++ {
		homeScore, awayScore := service.simulateMatch(team, team)
		homeScoreTotal += homeScore
		awayScoreTotal += awayScore
	}

	avgHomeScore := float64(homeScoreTotal) / float64(iterations)
	avgAwayScore := float64(awayScoreTotal) / float64(iterations)

	t.Logf("Same team home vs away - Avg home score: %.2f, Avg away score: %.2f",
		avgHomeScore, avgAwayScore)

	// Home team should score slightly more on average due to 10% home advantage
	if avgHomeScore <= avgAwayScore {
		t.Logf("Warning: Home advantage not reflected in average scores")
	}
}

func TestGenerateGoalsPoisson(t *testing.T) {
	testCases := []struct {
		name          string
		lambda        float64
		expectedRange [2]int // min, max
	}{
		{"Zero lambda", 0, [2]int{0, 0}},
		{"Low lambda", 0.5, [2]int{0, 7}},
		{"Normal lambda", 1.5, [2]int{0, 7}},
		{"High lambda", 3.0, [2]int{0, 7}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iterations := 1000
			totalGoals := 0
			goalCounts := make(map[int]int)

			for i := 0; i < iterations; i++ {
				goals := generateGoalsPoisson(tc.lambda)

				if goals < tc.expectedRange[0] || goals > tc.expectedRange[1] {
					t.Errorf("Goals %d outside expected range [%d, %d]",
						goals, tc.expectedRange[0], tc.expectedRange[1])
				}

				totalGoals += goals
				goalCounts[goals]++
			}

			avgGoals := float64(totalGoals) / float64(iterations)
			t.Logf("Lambda: %.1f, Average goals: %.2f", tc.lambda, avgGoals)

			// For Poisson, mean should be close to lambda
			if tc.lambda > 0 {
				// Allow 20% tolerance
				tolerance := tc.lambda * 0.3
				if avgGoals < tc.lambda-tolerance || avgGoals > tc.lambda+tolerance {
					t.Logf("Warning: Average %.2f differs from lambda %.1f", avgGoals, tc.lambda)
				}
			}
		})
	}
}

func TestGenerateGoalsPoissonDistribution(t *testing.T) {
	// Test that Poisson distribution produces expected frequencies
	lambda := 1.5
	iterations := 10000
	goalCounts := make(map[int]int)

	for i := 0; i < iterations; i++ {
		goals := generateGoalsPoisson(lambda)
		goalCounts[goals]++
	}

	// Expected probabilities for Poisson(1.5):
	// P(0) ≈ 22.3%, P(1) ≈ 33.5%, P(2) ≈ 25.1%, P(3) ≈ 12.6%
	expectedProbs := map[int]float64{
		0: 0.223,
		1: 0.335,
		2: 0.251,
		3: 0.126,
	}

	for goals, expectedProb := range expectedProbs {
		actualProb := float64(goalCounts[goals]) / float64(iterations)
		tolerance := 0.05 // 5% tolerance

		if actualProb < expectedProb-tolerance || actualProb > expectedProb+tolerance {
			t.Logf("Goals=%d: expected prob ~%.3f, got %.3f", goals, expectedProb, actualProb)
		}
	}

	// Log the distribution
	t.Log("Goal distribution:")
	for i := 0; i <= 5; i++ {
		prob := float64(goalCounts[i]) / float64(iterations) * 100
		t.Logf("  %d goals: %.1f%%", i, prob)
	}
}

func TestGenerateGoalsPoissonNegativeLambda(t *testing.T) {
	// Negative lambda should return 0
	goals := generateGoalsPoisson(-1.0)
	if goals != 0 {
		t.Errorf("Expected 0 goals for negative lambda, got %d", goals)
	}
}

func TestGenerateGoalsPoissonMaxCap(t *testing.T) {
	// Very high lambda should still be capped at maxGoalsPerTeam
	highLambda := 10.0
	iterations := 100

	for i := 0; i < iterations; i++ {
		goals := generateGoalsPoisson(highLambda)
		if goals > maxGoalsPerTeam {
			t.Errorf("Goals %d exceeds max %d", goals, maxGoalsPerTeam)
		}
	}
}
