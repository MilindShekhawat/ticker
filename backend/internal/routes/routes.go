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

	// Global priorities
	api.Get("/priorities", handlers.ListPriorities)
	api.Post("/priorities", handlers.CreatePriority)
	// Project priorities
	api.Get("/projects/:id/priorities", handlers.ListPriorities)
	api.Post("/projects/:id/priorities", handlers.CreatePriority)
	// Shared operations
	api.Patch("/priorities/:id", handlers.UpdatePriority)
	api.Patch("/priorities/:id/position", handlers.UpdatePriorityPosition)
	api.Delete("/priorities/:id", handlers.DeletePriority)

	// Global statuses
	api.Get("/statuses", handlers.ListStatuses)
	api.Post("/statuses", handlers.CreateStatus)
	// Project statuses
	api.Get("/projects/:id/statuses", handlers.ListStatuses)
	api.Post("/projects/:id/statuses", handlers.CreateStatus)
	// Shared operations
	api.Patch("/statuses/:id", handlers.UpdateStatus)
	api.Patch("/statuses/:id/position", handlers.UpdateStatusPosition)
	api.Delete("/statuses/:id", handlers.DeleteStatus)

	// Global tags
	api.Get("/tags", handlers.ListTags)
	api.Post("/tags", handlers.CreateTag)
	// Project tags
	api.Get("/projects/:id/tags", handlers.ListTags)
	api.Post("/projects/:id/tags", handlers.CreateTag)
	// Shared operations
	api.Patch("/tags/:id", handlers.UpdateTag)
	api.Delete("/tags/:id", handlers.DeleteTag)
	// Ticket tags
	api.Get("/tickets/:id/tags", handlers.ListTicketTags)
	api.Post("/tickets/:id/tags", handlers.AddTagToTicket)
	api.Delete("/tickets/:id/tags/:tag_id", handlers.RemoveTagFromTicket)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
