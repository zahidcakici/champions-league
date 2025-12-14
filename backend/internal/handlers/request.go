package handlers

type UpdateMatchResultRequest struct {
	HomeScore int `json:"homeScore" validate:"gte=0" example:"2"`
	AwayScore int `json:"awayScore" validate:"gte=0" example:"1"`
}

type CreateTeamRequest struct {
	Name  string `json:"name" validate:"required" example:"Team A"`
	Power int    `json:"power" validate:"gte=1,lte=100" example:"75"`
}

// Validate validates the request
func (r *UpdateMatchResultRequest) Validate() error {
	if r.HomeScore < 0 {
		return ErrInvalidHomeScore
	}
	if r.AwayScore < 0 {
		return ErrInvalidAwayScore
	}
	return nil
}
