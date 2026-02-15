package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Health check
	app.Get("/health", healthCheck)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
