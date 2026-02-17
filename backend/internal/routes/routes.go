package routes

import (
	"github.com/MilindShekhawat/ticker/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Health check
	app.Get("/health", healthCheck)

	// API v1 routes
	api := app.Group("/api/v1")

	// Projects
	api.Get("/projects", handlers.ListProjects)
	api.Post("/projects", handlers.CreateProject)
	api.Get("/projects/:id", handlers.GetProject)
	api.Patch("/projects/:id", handlers.UpdateProject)
	api.Delete("/projects/:id", handlers.DeleteProject)

	// Tickets
	api.Get("/tickets/:id", handlers.GetTicket)
	api.Get("/projects/:id/tickets", handlers.ListProjectTickets)
	api.Post("/projects/:id/tickets", handlers.CreateProjectTicket)
	api.Patch("/tickets/:id", handlers.UpdateTicket)
	api.Delete("/tickets/:id", handlers.DeleteTicket)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
