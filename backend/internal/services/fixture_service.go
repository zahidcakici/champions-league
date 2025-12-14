package services

import (
	"errors"

	"github.com/zahidcakici/champions-league/internal/models"
	"github.com/zahidcakici/champions-league/internal/repository"
)

type FixtureService interface {
	GenerateFixtures() ([]models.Match, error)
	GetAllFixtures() ([]models.Match, error)
	GetFixturesByWeek(week int) ([]models.Match, error)
}

type fixtureService struct {
	teamRepo   repository.TeamRepository
	matchRepo  repository.MatchRepository
	leagueRepo repository.LeagueStateRepository
}

func NewFixtureService(
	teamRepo repository.TeamRepository,
	matchRepo repository.MatchRepository,
	leagueRepo repository.LeagueStateRepository,
) FixtureService {
	return &fixtureService{
		teamRepo:   teamRepo,
		matchRepo:  matchRepo,
		leagueRepo: leagueRepo,
	}
}

func (s *fixtureService) GenerateFixtures() ([]models.Match, error) {
	// Check if fixtures already exist
	state, err := s.leagueRepo.Get()
	if err != nil {
		return nil, err
	}

	if state.FixturesCreated {
		return s.matchRepo.FindAll()
	}

	// Get all teams
	teams, err := s.teamRepo.FindAll()
	if err != nil {
		return nil, err
	}

	if len(teams) < 2 {
		return nil, errors.New("need at least 2 teams to generate fixtures")
	}

	// Generate round-robin fixtures (home and away)
	matches := s.generateRoundRobin(teams)

	// Save matches
	if err := s.matchRepo.CreateBatch(matches); err != nil {
		return nil, err
	}

	// Update league state
	state.FixturesCreated = true
	state.TotalWeeks = len(matches) / (len(teams) / 2)
	if err := s.leagueRepo.Update(state); err != nil {
		return nil, err
	}

	return s.matchRepo.FindAll()
}

// generateRoundRobin creates a round-robin schedule where each team plays every other team
// twice (home and away). Uses the circle method for fair scheduling.
func (s *fixtureService) generateRoundRobin(teams []models.Team) []models.Match {
	n := len(teams)

	// For n teams, we have n-1 rounds in single round-robin
	// For double round-robin (home and away), we have 2*(n-1) rounds
	// Each round has n/2 matches

	// First half: each team plays every other team once
	firstHalfMatches := s.generateSingleRoundRobin(teams)

	// Second half: reverse home/away
	secondHalf := make([]models.Match, len(firstHalfMatches))
	for i, m := range firstHalfMatches {
		secondHalf[i] = models.Match{
			Week:       m.Week + (n - 1),
			HomeTeamID: m.AwayTeamID,
			AwayTeamID: m.HomeTeamID,
			Played:     false,
		}
	}

	// Combine first and second half
	matches := make([]models.Match, 0, len(firstHalfMatches)+len(secondHalf))
	matches = append(matches, firstHalfMatches...)
	matches = append(matches, secondHalf...)

	return matches
}

// generateSingleRoundRobin creates matches where each team plays every other team once
func (s *fixtureService) generateSingleRoundRobin(teams []models.Team) []models.Match {
	n := len(teams)
	var matches []models.Match

	// Create a copy of team IDs for rotation
	teamIDs := make([]uint, n)
	for i, t := range teams {
		teamIDs[i] = t.ID
	}

	// Circle method: fix first team, rotate others
	for round := 0; round < n-1; round++ {
		week := round + 1

		for i := 0; i < n/2; i++ {
			home := i
			away := n - 1 - i

			// Alternate home/away for fairness
			if round%2 == 0 {
				matches = append(matches, models.Match{
					Week:       week,
					HomeTeamID: teamIDs[home],
					AwayTeamID: teamIDs[away],
					Played:     false,
				})
			} else {
				matches = append(matches, models.Match{
					Week:       week,
					HomeTeamID: teamIDs[away],
					AwayTeamID: teamIDs[home],
					Played:     false,
				})
			}
		}

		// Rotate: keep first element fixed, rotate the rest
		last := teamIDs[n-1]
		for i := n - 1; i > 1; i-- {
			teamIDs[i] = teamIDs[i-1]
		}
		teamIDs[1] = last
	}

	return matches
}

func (s *fixtureService) GetAllFixtures() ([]models.Match, error) {
	return s.matchRepo.FindAll()
}

func (s *fixtureService) GetFixturesByWeek(week int) ([]models.Match, error) {
	return s.matchRepo.FindByWeek(week)
}
