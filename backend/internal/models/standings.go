package models

// TeamStanding represents a team's position in the league table
type TeamStanding struct {
	TeamID         uint   `json:"team_id"`
	TeamName       string `json:"team_name"`
	Played         int    `json:"played"`
	Won            int    `json:"won"`
	Drawn          int    `json:"drawn"`
	Lost           int    `json:"lost"`
	GoalsFor       int    `json:"goals_for"`
	GoalsAgainst   int    `json:"goals_against"`
	GoalDifference int    `json:"goal_difference"`
	Points         int    `json:"points"`
}

// ChampionshipPrediction represents a team's probability of winning the championship
type ChampionshipPrediction struct {
	TeamID     uint    `json:"team_id"`
	TeamName   string  `json:"team_name"`
	Percentage float64 `json:"percentage"`
}

// SimulationState represents the complete state of the simulation
type SimulationState struct {
	LeagueState LeagueState              `json:"league_state"`
	Standings   []TeamStanding           `json:"standings"`
	CurrentWeek []MatchResult            `json:"current_week_results"`
	AllMatches  map[int][]MatchResult    `json:"all_matches"`
	Predictions []ChampionshipPrediction `json:"predictions"`
}
