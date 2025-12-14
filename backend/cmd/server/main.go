package main

import (
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/zahidcakici/champions-league/internal/config"
	"github.com/zahidcakici/champions-league/internal/database"
	"github.com/zahidcakici/champions-league/internal/handlers"
	"github.com/zahidcakici/champions-league/internal/handlers/docs"
	"github.com/zahidcakici/champions-league/internal/repository"
	"github.com/zahidcakici/champions-league/internal/routes"
	"github.com/zahidcakici/champions-league/internal/services"
)

// @title			Champions League Simulation API
// @version		1.0
// @description	A REST API for simulating a Champions League tournament
// @contact.name	API Support
// @contact.email	zahid.cakici@gmail.com
// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT
// @host		localhost:8080
// @BasePath	/api
// @schemes	http https
func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	leagueRepo := repository.NewLeagueStateRepository(db)

	// Initialize services
	teamService := services.NewTeamService(teamRepo)
	fixtureService := services.NewFixtureService(teamRepo, matchRepo, leagueRepo)
	simulationService := services.NewSimulationService(matchRepo, teamRepo, leagueRepo)
	standingsService := services.NewStandingsService(matchRepo, teamRepo, leagueRepo)

	// Initialize handlers
	teamHandler := handlers.NewTeamHandler(teamService)
	fixtureHandler := handlers.NewFixtureHandler(fixtureService)
	simulationHandler := handlers.NewSimulationHandler(simulationService, standingsService)
	standingsHandler := handlers.NewStandingsHandler(standingsService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	// Swagger documentation (embedded in binary)
	app.Use(swagger.New(swagger.Config{
		FileContent: docs.SwaggerJSON,
		Path:        "swagger",
		Title:       "Champions League API Documentation",
	}))

	// Setup routes
	routes.Setup(app, teamHandler, fixtureHandler, simulationHandler, standingsHandler)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
