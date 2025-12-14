package handlers

import "github.com/zahidcakici/champions-league/internal/models"

// teamToResponse converts a Team model to TeamResponse
func teamToResponse(team *models.Team) TeamResponse {
	return TeamResponse{
		ID:    team.ID,
		Name:  team.Name,
		Power: team.Power,
	}
}

// teamsToResponse converts a slice of Team models to TeamResponse slice
func teamsToResponse(teams []models.Team) []TeamResponse {
	responses := make([]TeamResponse, len(teams))
	for i := range teams {
		responses[i] = teamToResponse(&teams[i])
	}
	return responses
}

// matchToResponse converts a Match model to MatchResponse
func matchToResponse(match *models.Match) MatchResponse {
	return MatchResponse{
		ID:        match.ID,
		Week:      match.Week,
		HomeTeam:  teamToResponse(&match.HomeTeam),
		AwayTeam:  teamToResponse(&match.AwayTeam),
		HomeScore: match.HomeScore,
		AwayScore: match.AwayScore,
		Played:    match.Played,
	}
}

// matchesToResponse converts a slice of Match models to MatchResponse slice
func matchesToResponse(matches []models.Match) []MatchResponse {
	responses := make([]MatchResponse, len(matches))
	for i := range matches {
		responses[i] = matchToResponse(&matches[i])
	}
	return responses
}

// LeagueStateToResponse converts a LeagueState model to LeagueStateResponse
func LeagueStateToResponse(state *models.LeagueState) LeagueStateResponse {
	return LeagueStateResponse{
		CurrentWeek:     state.CurrentWeek,
		TotalWeeks:      state.TotalWeeks,
		FixturesCreated: state.FixturesCreated,
		Started:         state.Started,
		Completed:       state.Completed,
	}
}

// TeamStandingToResponse converts a TeamStanding model to TeamStandingResponse
func TeamStandingToResponse(standing *models.TeamStanding) TeamStandingResponse {
	return TeamStandingResponse{
		TeamID:         standing.TeamID,
		TeamName:       standing.TeamName,
		Played:         standing.Played,
		Won:            standing.Won,
		Drawn:          standing.Drawn,
		Lost:           standing.Lost,
		GoalsFor:       standing.GoalsFor,
		GoalsAgainst:   standing.GoalsAgainst,
		GoalDifference: standing.GoalDifference,
		Points:         standing.Points,
	}
}

// TeamStandingsToResponse converts a slice of TeamStanding models to TeamStandingResponse slice
func TeamStandingsToResponse(standings []models.TeamStanding) []TeamStandingResponse {
	responses := make([]TeamStandingResponse, len(standings))
	for i := range standings {
		responses[i] = TeamStandingToResponse(&standings[i])
	}
	return responses
}

// ChampionshipPredictionToResponse converts a ChampionshipPrediction model to ChampionshipPredictionResponse
func ChampionshipPredictionToResponse(prediction *models.ChampionshipPrediction) ChampionshipPredictionResponse {
	return ChampionshipPredictionResponse{
		TeamID:     prediction.TeamID,
		TeamName:   prediction.TeamName,
		Percentage: prediction.Percentage,
	}
}

// ChampionshipPredictionsToResponse converts a slice of ChampionshipPrediction models to response slice
func ChampionshipPredictionsToResponse(predictions []models.ChampionshipPrediction) []ChampionshipPredictionResponse {
	responses := make([]ChampionshipPredictionResponse, len(predictions))
	for i := range predictions {
		responses[i] = ChampionshipPredictionToResponse(&predictions[i])
	}
	return responses
}

// MatchResultToResponse converts a MatchResult model to MatchResultResponse
func MatchResultToResponse(result *models.MatchResult) MatchResultResponse {
	return MatchResultResponse{
		HomeTeamName: result.HomeTeamName,
		AwayTeamName: result.AwayTeamName,
		HomeScore:    result.HomeScore,
		AwayScore:    result.AwayScore,
	}
}

// MatchResultsToResponse converts a slice of MatchResult models to MatchResultResponse slice
func MatchResultsToResponse(results []models.MatchResult) []MatchResultResponse {
	responses := make([]MatchResultResponse, len(results))
	for i := range results {
		responses[i] = MatchResultToResponse(&results[i])
	}
	return responses
}

// AllMatchesToResponse converts a map of week to MatchResult slice to response format
func AllMatchesToResponse(allMatches map[int][]models.MatchResult) map[int][]MatchResultResponse {
	responses := make(map[int][]MatchResultResponse, len(allMatches))
	for week, results := range allMatches {
		responses[week] = MatchResultsToResponse(results)
	}
	return responses
}

// SimulationStateToResponse converts a SimulationState model to SimulationStateResponse
func SimulationStateToResponse(state *models.SimulationState) SimulationStateResponse {
	return SimulationStateResponse{
		LeagueState:        LeagueStateToResponse(&state.LeagueState),
		Standings:          TeamStandingsToResponse(state.Standings),
		CurrentWeekResults: MatchResultsToResponse(state.CurrentWeek),
		AllMatches:         AllMatchesToResponse(state.AllMatches),
		Predictions:        ChampionshipPredictionsToResponse(state.Predictions),
	}
}
