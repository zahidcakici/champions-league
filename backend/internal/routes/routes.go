package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zahidcakici/champions-league/internal/handlers"
)

func Setup(
	app *fiber.App,
	teamHandler *handlers.TeamHandler,
	fixtureHandler *handlers.FixtureHandler,
	simulationHandler *handlers.SimulationHandler,
	standingsHandler *handlers.StandingsHandler,
) {
	api := app.Group("/api")

	// Team routes
	teams := api.Group("/teams")
	teams.Get("/", teamHandler.GetAllTeams)
	teams.Post("/", teamHandler.CreateTeam)
	teams.Delete("/:id", teamHandler.DeleteTeam)

	// Fixture routes
	fixtures := api.Group("/fixtures")
	fixtures.Get("/", fixtureHandler.GetAllFixtures)
	fixtures.Get("/:week", fixtureHandler.GetFixturesByWeek)
	fixtures.Post("/generate", fixtureHandler.GenerateFixtures)

	// Simulation routes
	simulation := api.Group("/simulation")
	simulation.Get("/state", simulationHandler.GetState)
	simulation.Post("/play-week", simulationHandler.PlayNextWeek)
	simulation.Post("/play-all", simulationHandler.PlayAllWeeks)
	simulation.Put("/match/:id", simulationHandler.UpdateMatchResult)
	simulation.Post("/reset", simulationHandler.ResetSimulation)

	// Standings routes
	api.Get("/standings", standingsHandler.GetStandings)
	api.Get("/predictions", standingsHandler.GetPredictions)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
