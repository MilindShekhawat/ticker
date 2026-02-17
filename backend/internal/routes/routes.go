package routes

import (
	"github.com/MilindShekhawat/ticker/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Health check
	app.Get("/health", healthCheck)

	// API routes
	api := app.Group("/api")

	// Tickets
	api.Get("/projects/:id/tickets", handlers.GetTickets)
	api.Post("/projects/:id/tickets", handlers.CreateTicket)
	api.Get("/tickets/:id", handlers.GetTicket)
	api.Patch("/tickets/:id", handlers.UpdateTicket)
	api.Delete("/tickets/:id", handlers.DeleteTicket)
	// api.Patch("/projects/:id/tickets/bulk", handlers.BulkUpdateTickets)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
