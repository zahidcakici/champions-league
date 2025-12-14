package services

import (
	"errors"
	"math"
	"math/rand"

	"github.com/zahidcakici/champions-league/internal/models"
	"github.com/zahidcakici/champions-league/internal/repository"
)

const (
	homeAdvantageFactor = 1.1
	baseExpectedGoals   = 1.5
	maxGoalsPerTeam     = 7
)

type SimulationService interface {
	PlayNextWeek() ([]models.Match, error)
	PlayAllWeeks() (map[int][]models.Match, error)
	UpdateMatchResult(matchID uint, homeScore, awayScore int) error
	ResetSimulation() error
	GetCurrentState() (*models.SimulationState, error)
}

type simulationService struct {
	matchRepo  repository.MatchRepository
	teamRepo   repository.TeamRepository
	leagueRepo repository.LeagueStateRepository
}

func NewSimulationService(
	matchRepo repository.MatchRepository,
	teamRepo repository.TeamRepository,
	leagueRepo repository.LeagueStateRepository,
) SimulationService {
	return &simulationService{
		matchRepo:  matchRepo,
		teamRepo:   teamRepo,
		leagueRepo: leagueRepo,
	}
}

func (s *simulationService) PlayNextWeek() ([]models.Match, error) {
	state, err := s.leagueRepo.Get()
	if err != nil {
		return nil, err
	}

	if !state.FixturesCreated {
		return nil, errors.New("fixtures not generated yet")
	}

	if state.Completed {
		return nil, errors.New("league already completed")
	}

	nextWeek := state.CurrentWeek + 1
	matches, err := s.matchRepo.FindByWeek(nextWeek)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, errors.New("no matches found for this week")
	}

	// Simulate each match
	for i := range matches {
		if !matches[i].Played {
			homeScore, awayScore := s.simulateMatch(&matches[i].HomeTeam, &matches[i].AwayTeam)
			matches[i].HomeScore = &homeScore
			matches[i].AwayScore = &awayScore
			matches[i].Played = true
			if err := s.matchRepo.Update(&matches[i]); err != nil {
				return nil, err
			}
		}
	}

	// Update league state
	state.CurrentWeek = nextWeek
	state.Started = true
	if nextWeek >= state.TotalWeeks {
		state.Completed = true
	}
	if err := s.leagueRepo.Update(state); err != nil {
		return nil, err
	}

	return matches, nil
}

func (s *simulationService) PlayAllWeeks() (map[int][]models.Match, error) {
	state, err := s.leagueRepo.Get()
	if err != nil {
		return nil, err
	}

	if !state.FixturesCreated {
		return nil, errors.New("fixtures not generated yet")
	}

	results := make(map[int][]models.Match)

	for !state.Completed {
		matches, err := s.PlayNextWeek()
		if err != nil {
			return nil, err
		}

		// Refresh state
		state, err = s.leagueRepo.Get()
		if err != nil {
			return nil, err
		}

		results[state.CurrentWeek] = matches
	}

	return results, nil
}

// simulateMatch generates a match result based on team powers
// Uses weighted random algorithm with home advantage
// Each team's expected goals depends on their power relative to opponent's power
func (s *simulationService) simulateMatch(homeTeam, awayTeam *models.Team) (int, int) {
	// Calculate effective powers
	homePower := float64(homeTeam.Power) * homeAdvantageFactor
	awayPower := float64(awayTeam.Power)

	// Total power for relative calculations
	totalPower := homePower + awayPower

	// Generate goals based on relative power
	// Expected goals = base * (own power relative to total)
	// This means stronger opponents reduce your expected goals
	homeExpectedGoals := baseExpectedGoals * 2 * (homePower / totalPower)
	awayExpectedGoals := baseExpectedGoals * 2 * (awayPower / totalPower)

	homeGoals := generateGoalsPoisson(homeExpectedGoals)
	awayGoals := generateGoalsPoisson(awayExpectedGoals)

	return homeGoals, awayGoals
}

// generateGoals generates a realistic goal count using a simplified Poisson-like distribution - Old method
func (s *simulationService) generateGoals(expectedGoals float64) int {
	goals := 0

	// Simulate multiple "chances" to score
	chances := 10
	scoreProbability := expectedGoals / float64(chances)

	for i := 0; i < chances; i++ {
		if rand.Float64() < scoreProbability {
			goals++
		}
	}

	// Cap at reasonable maximum
	if goals > maxGoalsPerTeam {
		goals = maxGoalsPerTeam
	}

	return goals
}

// generateGoalsPoisson uses Poisson distribution for realistic goal counts
// This models the discrete nature of goals where:
// - 0-1 goals are most common
// - 2-3 goals are fairly common
// - 4+ goals are rare but possible
func generateGoalsPoisson(lambda float64) int {
	if lambda <= 0 {
		return 0
	}

	// Poisson distribution using inverse transform sampling
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}

	return min(k-1, maxGoalsPerTeam)
}

func (s *simulationService) UpdateMatchResult(matchID uint, homeScore, awayScore int) error {
	match, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return err
	}

	match.HomeScore = &homeScore
	match.AwayScore = &awayScore
	match.Played = true

	return s.matchRepo.Update(match)
}

func (s *simulationService) ResetSimulation() error {
	// Delete all matches
	if err := s.matchRepo.DeleteAll(); err != nil {
		return err
	}

	// Reset league state
	if err := s.leagueRepo.Reset(); err != nil {
		return err
	}

	return nil
}

func (s *simulationService) GetCurrentState() (*models.SimulationState, error) {
	leagueState, err := s.leagueRepo.Get()
	if err != nil {
		return nil, err
	}

	return &models.SimulationState{
		LeagueState: *leagueState,
	}, nil
}
