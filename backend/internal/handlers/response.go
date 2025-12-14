package handlers

// APIResponse is the standard API response wrapper
// @Description Standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}

// APIErrorResponse is the standard API error response
// @Description Standard API error response
type APIErrorResponse struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Something went wrong"`
}

// TeamResponse represents a team in API responses
// @Description Team information
type TeamResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Manchester City"`
	Power int    `json:"power" example:"90"`
}

// MatchResponse represents a match in API responses
// @Description Match information
type MatchResponse struct {
	ID        uint         `json:"id" example:"1"`
	Week      int          `json:"week" example:"1"`
	HomeTeam  TeamResponse `json:"homeTeam"`
	AwayTeam  TeamResponse `json:"awayTeam"`
	HomeScore *int         `json:"homeScore" example:"2"`
	AwayScore *int         `json:"awayScore" example:"1"`
	Played    bool         `json:"played" example:"true"`
}

// LeagueStateResponse represents the league state in API responses
// @Description Current league state
type LeagueStateResponse struct {
	CurrentWeek     int  `json:"currentWeek" example:"3"`
	TotalWeeks      int  `json:"totalWeeks" example:"6"`
	FixturesCreated bool `json:"fixturesCreated" example:"true"`
	Started         bool `json:"started" example:"true"`
	Completed       bool `json:"completed" example:"false"`
}

// TeamStandingResponse represents a team's standing in the league table
// @Description Team standing in league table
type TeamStandingResponse struct {
	TeamID         uint   `json:"teamId" example:"1"`
	TeamName       string `json:"teamName" example:"Manchester City"`
	Played         int    `json:"played" example:"3"`
	Won            int    `json:"won" example:"2"`
	Drawn          int    `json:"drawn" example:"1"`
	Lost           int    `json:"lost" example:"0"`
	GoalsFor       int    `json:"goalsFor" example:"7"`
	GoalsAgainst   int    `json:"goalsAgainst" example:"3"`
	GoalDifference int    `json:"goalDifference" example:"4"`
	Points         int    `json:"points" example:"7"`
}

// ChampionshipPredictionResponse represents a team's championship probability
// @Description Championship prediction for a team
type ChampionshipPredictionResponse struct {
	TeamID     uint    `json:"teamId" example:"1"`
	TeamName   string  `json:"teamName" example:"Manchester City"`
	Percentage float64 `json:"percentage" example:"45.5"`
}

// MatchResultResponse represents a played match result
// @Description Match result
type MatchResultResponse struct {
	HomeTeamName string `json:"homeTeamName" example:"Chelsea"`
	AwayTeamName string `json:"awayTeamName" example:"Arsenal"`
	HomeScore    int    `json:"homeScore" example:"2"`
	AwayScore    int    `json:"awayScore" example:"1"`
}

// SimulationStateResponse represents the complete simulation state
// @Description Complete simulation state including standings and predictions
type SimulationStateResponse struct {
	LeagueState        LeagueStateResponse              `json:"leagueState"`
	Standings          []TeamStandingResponse           `json:"standings"`
	CurrentWeekResults []MatchResultResponse            `json:"currentWeekResults"`
	AllMatches         map[int][]MatchResultResponse    `json:"allMatches"`
	Predictions        []ChampionshipPredictionResponse `json:"predictions"`
}

// TeamsListResponse is the response for GET /teams
// @Description List of all teams
type TeamsListResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    []TeamResponse `json:"data"`
}

// FixturesListResponse is the response for GET /fixtures
// @Description List of all fixtures
type FixturesListResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    []MatchResponse `json:"data"`
}

// StandingsListResponse is the response for GET /standings
// @Description Current league standings
type StandingsListResponse struct {
	Success bool                   `json:"success" example:"true"`
	Data    []TeamStandingResponse `json:"data"`
}

// PredictionsListResponse is the response for GET /predictions
// @Description Championship predictions
type PredictionsListResponse struct {
	Success bool                             `json:"success" example:"true"`
	Data    []ChampionshipPredictionResponse `json:"data"`
}

// SimulationStateFullResponse is the response for simulation state endpoints
// @Description Full simulation state response
type SimulationStateFullResponse struct {
	Success bool                    `json:"success" example:"true"`
	Data    SimulationStateResponse `json:"data"`
}

// MessageResponse is a simple message response
// @Description Simple message response
type MessageResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    MessageData `json:"data"`
}

// MessageData contains a simple message
type MessageData struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
