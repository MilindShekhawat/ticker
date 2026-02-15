package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/MilindShekhawat/ticker/internal/handlers"
)

func Routes(app *fiber.App) {
	// Health check
	app.Get("/health", healthCheck)

	// API routes
	api := app.Group("/api")

	// Tickets
	api.Get("/tickets", handlers.GetTickets)
	api.Get("/tickets/:id", handlers.GetTicket)
	api.Post("/tickets", handlers.CreateTicket)
	api.Put("/tickets/:id", handlers.UpdateTicket)
	api.Delete("/tickets/:id", handlers.DeleteTicket)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
