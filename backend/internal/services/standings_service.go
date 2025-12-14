package services

import (
	"math"
	"sort"

	"github.com/zahidcakici/champions-league/internal/models"
	"github.com/zahidcakici/champions-league/internal/repository"
)

type StandingsService interface {
	GetStandings() ([]models.TeamStanding, error)
	GetPredictions() ([]models.ChampionshipPrediction, error)
	GetFullState() (*models.SimulationState, error)
}

type standingsService struct {
	matchRepo  repository.MatchRepository
	teamRepo   repository.TeamRepository
	leagueRepo repository.LeagueStateRepository
}

func NewStandingsService(
	matchRepo repository.MatchRepository,
	teamRepo repository.TeamRepository,
	leagueRepo repository.LeagueStateRepository,
) StandingsService {
	return &standingsService{
		matchRepo:  matchRepo,
		teamRepo:   teamRepo,
		leagueRepo: leagueRepo,
	}
}

func (s *standingsService) GetStandings() ([]models.TeamStanding, error) {
	teams, err := s.teamRepo.FindAll()
	if err != nil {
		return nil, err
	}

	matches, err := s.matchRepo.FindPlayedMatches()
	if err != nil {
		return nil, err
	}

	// Initialize standings for all teams
	standingsMap := make(map[uint]*models.TeamStanding)
	for _, team := range teams {
		standingsMap[team.ID] = &models.TeamStanding{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	// Calculate standings from played matches
	for _, match := range matches {
		if match.HomeScore == nil || match.AwayScore == nil {
			continue
		}

		homeStanding := standingsMap[match.HomeTeamID]
		awayStanding := standingsMap[match.AwayTeamID]

		homeStanding.Played++
		awayStanding.Played++

		homeStanding.GoalsFor += *match.HomeScore
		homeStanding.GoalsAgainst += *match.AwayScore
		awayStanding.GoalsFor += *match.AwayScore
		awayStanding.GoalsAgainst += *match.HomeScore

		switch {
		case *match.HomeScore > *match.AwayScore:
			// Home win
			homeStanding.Won++
			homeStanding.Points += 3
			awayStanding.Lost++
		case *match.HomeScore < *match.AwayScore:
			// Away win
			awayStanding.Won++
			awayStanding.Points += 3
			homeStanding.Lost++
		default:
			// Draw
			homeStanding.Drawn++
			awayStanding.Drawn++
			homeStanding.Points++
			awayStanding.Points++
		}
	}

	// Calculate goal difference and convert to slice
	var standings []models.TeamStanding
	for _, standing := range standingsMap {
		standing.GoalDifference = standing.GoalsFor - standing.GoalsAgainst
		standings = append(standings, *standing)
	}

	// Sort by points (desc), then goal difference (desc), then goals for (desc)
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		if standings[i].GoalDifference != standings[j].GoalDifference {
			return standings[i].GoalDifference > standings[j].GoalDifference
		}
		return standings[i].GoalsFor > standings[j].GoalsFor
	})

	return standings, nil
}

func (s *standingsService) GetPredictions() ([]models.ChampionshipPrediction, error) {
	state, err := s.leagueRepo.Get()
	if err != nil {
		return nil, err
	}

	standings, err := s.GetStandings()
	if err != nil {
		return nil, err
	}

	predictions := make([]models.ChampionshipPrediction, len(standings))

	// Predictions only start when there are 3 or fewer weeks remaining
	// This gives enough data for meaningful predictions while keeping excitement
	remainingWeeks := state.TotalWeeks - state.CurrentWeek
	if remainingWeeks > 3 {
		for i, standing := range standings {
			predictions[i] = models.ChampionshipPrediction{
				TeamID:     standing.TeamID,
				TeamName:   standing.TeamName,
				Percentage: 0,
			}
		}
		return predictions, nil
	}

	// Calculate max remaining points (each team plays 1 match per week)
	maxRemainingPoints := remainingWeeks * 3

	// Get the leader's points
	if len(standings) == 0 {
		return predictions, nil
	}

	leaderPoints := standings[0].Points

	// Calculate championship probability based on points gap and remaining matches
	totalWeight := 0.0
	weights := make([]float64, len(standings))

	for i, standing := range standings {
		pointsGap := leaderPoints - standing.Points

		// If a team can't mathematically catch up, their chance is 0
		if pointsGap > maxRemainingPoints {
			weights[i] = 0
		} else {
			// Weight based on current points and ability to catch up
			// Teams closer to the leader have higher probability
			// Also factor in goal difference as tiebreaker potential
			weight := float64(standing.Points+1) * math.Pow(0.7, float64(pointsGap))

			// Bonus for positive goal difference
			if standing.GoalDifference > 0 {
				weight *= 1.1
			}

			weights[i] = weight
		}
		totalWeight += weights[i]
	}

	// Normalize to percentages
	for i, standing := range standings {
		percentage := 0.0
		if totalWeight > 0 {
			percentage = (weights[i] / totalWeight) * 100
		}

		// Round to nearest integer
		predictions[i] = models.ChampionshipPrediction{
			TeamID:     standing.TeamID,
			TeamName:   standing.TeamName,
			Percentage: math.Round(percentage),
		}
	}

	// Ensure percentages sum to 100%
	s.normalizePercentages(predictions)

	return predictions, nil
}

func (s *standingsService) normalizePercentages(predictions []models.ChampionshipPrediction) {
	total := 0.0
	for _, p := range predictions {
		total += p.Percentage
	}

	if total == 0 {
		// Equal distribution if no predictions
		equal := 100.0 / float64(len(predictions))
		for i := range predictions {
			predictions[i].Percentage = math.Round(equal)
		}
		return
	}

	// Adjust to sum to 100
	if total != 100 {
		diff := 100 - total
		// Add/subtract difference from the leader
		predictions[0].Percentage += diff
	}
}

func (s *standingsService) GetFullState() (*models.SimulationState, error) {
	leagueState, err := s.leagueRepo.Get()
	if err != nil {
		return nil, err
	}

	standings, err := s.GetStandings()
	if err != nil {
		return nil, err
	}

	predictions, err := s.GetPredictions()
	if err != nil {
		return nil, err
	}

	// Get current week results
	var currentWeekResults []models.MatchResult
	if leagueState.CurrentWeek > 0 {
		weekMatches, matchErr := s.matchRepo.FindByWeek(leagueState.CurrentWeek)
		if matchErr != nil {
			return nil, matchErr
		}
		for _, m := range weekMatches {
			if m.Played && m.HomeScore != nil && m.AwayScore != nil {
				currentWeekResults = append(currentWeekResults, models.MatchResult{
					HomeTeamName: m.HomeTeam.Name,
					AwayTeamName: m.AwayTeam.Name,
					HomeScore:    *m.HomeScore,
					AwayScore:    *m.AwayScore,
				})
			}
		}
	}

	// Get all matches grouped by week
	allMatches := make(map[int][]models.MatchResult)
	matches, err := s.matchRepo.FindAll()
	if err != nil {
		return nil, err
	}
	for _, m := range matches {
		result := models.MatchResult{
			HomeTeamName: m.HomeTeam.Name,
			AwayTeamName: m.AwayTeam.Name,
		}
		if m.Played && m.HomeScore != nil && m.AwayScore != nil {
			result.HomeScore = *m.HomeScore
			result.AwayScore = *m.AwayScore
		}
		allMatches[m.Week] = append(allMatches[m.Week], result)
	}

	return &models.SimulationState{
		LeagueState: *leagueState,
		Standings:   standings,
		CurrentWeek: currentWeekResults,
		AllMatches:  allMatches,
		Predictions: predictions,
	}, nil
}
