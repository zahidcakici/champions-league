package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zahidcakici/champions-league/internal/models"
	"github.com/zahidcakici/champions-league/internal/services"
)

type TeamHandler struct {
	teamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// GetAllTeams returns all teams in the league
//
//	@Summary		Get all teams
//	@Description	Returns all teams participating in the tournament with their power ratings
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	TeamsListResponse	"Success response with teams array"
//	@Failure		500	{object}	APIErrorResponse	"Internal server error"
//	@Router			/teams [get]
func (h *TeamHandler) GetAllTeams(c *fiber.Ctx) error {
	teams, err := h.teamService.GetAllTeams()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, teamsToResponse(teams))
}

// CreateTeam creates a new team
//
//	@Summary		Create a new team
//	@Description	Creates a new team with a specified name and power rating
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			team	body		CreateTeamRequest	true	"Team creation payload"
//	@Success		201		{object}	TeamResponse			"Success response with created team"
//	@Failure		400		{object}	APIErrorResponse		"Bad request (e.g., invalid input)"
//	@Failure		500		{object}	APIErrorResponse		"Internal server error"
//	@Router			/teams [post]
func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	var req CreateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validation
	if req.Name == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Team name is required")
	}
	if req.Power < 1 || req.Power > 100 {
		return ErrorResponse(c, fiber.StatusBadRequest, "Team power must be between 1 and 100")
	}

	err := h.teamService.CreateTeam(req.Name, req.Power)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, teamToResponse(&models.Team{
		Name:  req.Name,
		Power: req.Power,
	}))
}

// DeleteTeam deletes a team by ID
//
//	@Summary		Delete a team
//	@Description	Deletes a team by its ID. Can only be done before fixtures are generated.
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Team ID"
//	@Success		200	{object}	APISuccessResponse	"Success response"
//	@Failure		400	{object}	APIErrorResponse	"Bad request (e.g., invalid ID)"
//	@Failure		500	{object}	APIErrorResponse	"Internal server error"
//	@Router			/teams/{id} [delete]
func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid team ID")
	}

	if err := h.teamService.DeleteTeam(uint(id)); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, fiber.Map{"deleted": true})
}
