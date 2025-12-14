package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zahidcakici/champions-league/internal/services"
)

type StandingsHandler struct {
	standingsService services.StandingsService
}

func NewStandingsHandler(standingsService services.StandingsService) *StandingsHandler {
	return &StandingsHandler{standingsService: standingsService}
}

// GetStandings returns the current league standings
//
//	@Summary		Get league standings
//	@Description	Returns the current league table with points, goals, and positions
//	@Tags			Standings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	StandingsListResponse	"Success response with standings array"
//	@Failure		500	{object}	APIErrorResponse		"Internal server error"
//	@Router			/standings [get]
func (h *StandingsHandler) GetStandings(c *fiber.Ctx) error {
	standings, err := h.standingsService.GetStandings()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, TeamStandingsToResponse(standings))
}

// GetPredictions returns championship predictions
//
//	@Summary		Get championship predictions
//	@Description	Returns the probability of each team winning the championship (available from week 4)
//	@Tags			Standings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	PredictionsListResponse	"Success response with predictions array"
//	@Failure		500	{object}	APIErrorResponse		"Internal server error"
//	@Router			/predictions [get]
func (h *StandingsHandler) GetPredictions(c *fiber.Ctx) error {
	predictions, err := h.standingsService.GetPredictions()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, ChampionshipPredictionsToResponse(predictions))
}
