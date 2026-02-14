package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/MilindShekhawat/ticker/internal/db"
)

func main() {
	// db connection
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// migrations
	cwd, err := os.Getwd()
	if err != nil {
	    log.Fatalf("Failed to get working directory: %v", err)
	}
	migrationsPath := filepath.Join(cwd, "migrations")

	if err := db.RunMigrations(migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// http server
	app := fiber.New(fiber.Config{
		AppName: "Ticker v0.1.0",
	})

	// middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// API routes ----------------------------------
	api := app.Group("/api")
	api.Get("/tickets", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "List tickets - coming soon"})
	})

	// start server
	log.Println("Server starting on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
