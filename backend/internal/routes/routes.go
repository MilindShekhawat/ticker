package routes

import (
	"github.com/MilindShekhawat/ticker/internal/handlers"
	"github.com/MilindShekhawat/ticker/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Health check
	app.Get("/health", healthCheck)

	// API v1 routes
	api := app.Group("/api/v1")
	protected := api.Group("", middleware.AuthRequired)

	// Public Auth
	api.Post("/auth/signup", handlers.Signup)
	api.Post("/auth/login", handlers.Login)
	// Protected Auth
	protected.Post("/auth/logout", handlers.Logout)
	protected.Get("/auth/me", handlers.GetCurrentUser)

	// Projects
	protected.Get("/projects", handlers.ListProjects)
	protected.Post("/projects", handlers.CreateProject)
	protected.Get("/projects/:id", handlers.GetProject)
	protected.Patch("/projects/:id", handlers.UpdateProject)
	protected.Delete("/projects/:id", handlers.DeleteProject)

	// Tickets
	protected.Get("/projects/:id/tickets", handlers.ListProjectTickets)
	protected.Post("/projects/:id/tickets", handlers.CreateProjectTicket)
	protected.Get("/tickets/:id", handlers.GetTicket)
	protected.Patch("/tickets/:id", handlers.UpdateTicket)
	protected.Delete("/tickets/:id", handlers.DeleteTicket)

	// Global priorities
	protected.Get("/priorities", handlers.ListPriorities)
	protected.Post("/priorities", handlers.CreatePriority)
	// Project priorities
	protected.Get("/projects/:id/priorities", handlers.ListPriorities)
	protected.Post("/projects/:id/priorities", handlers.CreatePriority)
	// Shared operations
	protected.Patch("/priorities/:id", handlers.UpdatePriority)
	protected.Patch("/priorities/:id/position", handlers.UpdatePriorityPosition)
	protected.Delete("/priorities/:id", handlers.DeletePriority)

	// Global statuses
	protected.Get("/statuses", handlers.ListStatuses)
	protected.Post("/statuses", handlers.CreateStatus)
	// Project statuses
	protected.Get("/projects/:id/statuses", handlers.ListStatuses)
	protected.Post("/projects/:id/statuses", handlers.CreateStatus)
	// Shared operations
	protected.Patch("/statuses/:id", handlers.UpdateStatus)
	protected.Patch("/statuses/:id/position", handlers.UpdateStatusPosition)
	protected.Delete("/statuses/:id", handlers.DeleteStatus)

	// Global tags
	protected.Get("/tags", handlers.ListTags)
	protected.Post("/tags", handlers.CreateTag)
	// Project tags
	protected.Get("/projects/:id/tags", handlers.ListTags)
	protected.Post("/projects/:id/tags", handlers.CreateTag)
	// Shared operations
	protected.Patch("/tags/:id", handlers.UpdateTag)
	protected.Delete("/tags/:id", handlers.DeleteTag)

	// Comments
	protected.Get("/tickets/:id/comments", handlers.ListComments)
	protected.Post("/tickets/:id/comments", handlers.CreateComment)
	protected.Patch("/comments/:id", handlers.UpdateComment)
	protected.Delete("/comments/:id", handlers.DeleteComment)

	// Activity
	protected.Get("/tickets/:id/activity", handlers.GetTicketActivity)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
