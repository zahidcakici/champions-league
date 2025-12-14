package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zahidcakici/champions-league/internal/services"
)

type FixtureHandler struct {
	fixtureService services.FixtureService
}

func NewFixtureHandler(fixtureService services.FixtureService) *FixtureHandler {
	return &FixtureHandler{fixtureService: fixtureService}
}

// GenerateFixtures creates the fixture schedule
//
//	@Summary		Generate fixtures
//	@Description	Creates a round-robin fixture schedule for all teams (6 weeks, home and away)
//	@Tags			Fixtures
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	FixturesListResponse	"Success response with generated fixtures"
//	@Failure		500	{object}	APIErrorResponse		"Internal server error"
//	@Router			/fixtures/generate [post]
func (h *FixtureHandler) GenerateFixtures(c *fiber.Ctx) error {
	fixtures, err := h.fixtureService.GenerateFixtures()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, matchesToResponse(fixtures))
}

// GetAllFixtures returns all fixtures
//
//	@Summary		Get all fixtures
//	@Description	Returns all fixtures across all weeks
//	@Tags			Fixtures
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	FixturesListResponse	"Success response with fixtures array"
//	@Failure		500	{object}	APIErrorResponse		"Internal server error"
//	@Router			/fixtures [get]
func (h *FixtureHandler) GetAllFixtures(c *fiber.Ctx) error {
	fixtures, err := h.fixtureService.GetAllFixtures()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, matchesToResponse(fixtures))
}

// GetFixturesByWeek returns fixtures for a specific week
//
//	@Summary		Get fixtures by week
//	@Description	Returns all fixtures for a specific week number
//	@Tags			Fixtures
//	@Accept			json
//	@Produce		json
//	@Param			week	path		int						true	"Week number (1-6)"
//	@Success		200		{object}	FixturesListResponse	"Success response with fixtures for the week"
//	@Failure		400		{object}	APIErrorResponse		"Invalid week number"
//	@Failure		500		{object}	APIErrorResponse		"Internal server error"
//	@Router			/fixtures/{week} [get]
func (h *FixtureHandler) GetFixturesByWeek(c *fiber.Ctx) error {
	week, err := c.ParamsInt("week")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid week number")
	}

	fixtures, err := h.fixtureService.GetFixturesByWeek(week)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, matchesToResponse(fixtures))
}
