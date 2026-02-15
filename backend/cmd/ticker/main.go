package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/MilindShekhawat/ticker/internal/config"
	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/routes"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := db.Init(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	migrationsPath := filepath.Join(cwd, "migrations")

	if err := db.RunMigrations(migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Ticker v0.1.0",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] | ${latency} | ${status} - ${method} ${path}\n",
	}))
	app.Use(recover.New())

	corsConfig := cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}

	app.Use(cors.New(corsConfig))

	// Api routes
	routes.Routes(app)

	// Start server
	log.Printf("Environment: %s\n", cfg.Env)
	log.Printf("Server ready on http://%s\n", cfg.Address())
	if cfg.Host == "0.0.0.0" {
		log.Printf("Network access enabled (accessible from other devices)\n")
	}

	if err := app.Listen(cfg.Address()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
