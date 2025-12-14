package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zahidcakici/champions-league/internal/services"
)

type SimulationHandler struct {
	simulationService services.SimulationService
	standingsService  services.StandingsService
}

func NewSimulationHandler(
	simulationService services.SimulationService,
	standingsService services.StandingsService,
) *SimulationHandler {
	return &SimulationHandler{
		simulationService: simulationService,
		standingsService:  standingsService,
	}
}

// PlayNextWeek simulates the next week of matches
//
//	@Summary		Play next week
//	@Description	Simulates all matches for the next week and returns updated state
//	@Tags			Simulation
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SimulationStateFullResponse	"Success response with updated simulation state"
//	@Failure		400	{object}	APIErrorResponse			"Bad request (e.g., no more weeks to play)"
//	@Failure		500	{object}	APIErrorResponse			"Internal server error"
//	@Router			/simulation/play-week [post]
func (h *SimulationHandler) PlayNextWeek(c *fiber.Ctx) error {
	_, err := h.simulationService.PlayNextWeek()
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Return full state after playing
	state, err := h.standingsService.GetFullState()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, SimulationStateToResponse(state))
}

// PlayAllWeeks simulates all remaining weeks
//
//	@Summary		Play all remaining weeks
//	@Description	Simulates all remaining matches until the season is complete
//	@Tags			Simulation
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SimulationStateFullResponse	"Success response with final simulation state"
//	@Failure		400	{object}	APIErrorResponse			"Bad request (e.g., season already complete)"
//	@Failure		500	{object}	APIErrorResponse			"Internal server error"
//	@Router			/simulation/play-all [post]
func (h *SimulationHandler) PlayAllWeeks(c *fiber.Ctx) error {
	_, err := h.simulationService.PlayAllWeeks()
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Return full state after playing all
	state, err := h.standingsService.GetFullState()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, SimulationStateToResponse(state))
}

// UpdateMatchResult updates the result of a specific match
//
//	@Summary		Update match result
//	@Description	Manually update the score of a specific match
//	@Tags			Simulation
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Match ID"
//	@Param			body	body		UpdateMatchResultRequest	true	"Match score update"
//	@Success		200		{object}	SimulationStateFullResponse	"Success response with updated simulation state"
//	@Failure		400		{object}	APIErrorResponse			"Invalid match ID or request body"
//	@Failure		500		{object}	APIErrorResponse			"Internal server error"
//	@Router			/simulation/match/{id} [put]
func (h *SimulationHandler) UpdateMatchResult(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid match ID")
	}

	var req UpdateMatchResultRequest
	if err = c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err = req.Validate(); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err = h.simulationService.UpdateMatchResult(uint(id), req.HomeScore, req.AwayScore); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Return full state after update
	state, err := h.standingsService.GetFullState()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, SimulationStateToResponse(state))
}

// ResetSimulation resets the simulation to initial state (keeps fixtures)
//
//	@Summary		Reset simulation
//	@Description	Resets all match results and league state while keeping fixtures
//	@Tags			Simulation
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	MessageResponse		"Success response with reset confirmation"
//	@Failure		500	{object}	APIErrorResponse	"Internal server error"
//	@Router			/simulation/reset [post]
func (h *SimulationHandler) ResetSimulation(c *fiber.Ctx) error {
	if err := h.simulationService.ResetSimulation(); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, MessageData{Message: "Simulation reset successfully"})
}

// GetState returns the current simulation state
//
//	@Summary		Get simulation state
//	@Description	Returns the complete current state including standings, predictions, and match results
//	@Tags			Simulation
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SimulationStateFullResponse	"Success response with full simulation state"
//	@Failure		500	{object}	APIErrorResponse			"Internal server error"
//	@Router			/simulation/state [get]
func (h *SimulationHandler) GetState(c *fiber.Ctx) error {
	state, err := h.standingsService.GetFullState()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, SimulationStateToResponse(state))
}
