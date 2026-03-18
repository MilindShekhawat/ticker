package routes

import (
	"github.com/MilindShekhawat/ticker/internal/handlers"
	"github.com/MilindShekhawat/ticker/internal/middleware"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	// Health check
	app.Get("/health", healthCheck)

	// API v1 routes
	api := app.Group("/api/v1")

	// Auth
	userStore := store.NewUserStore()
	sessionStore := store.NewSessionStore()
	authService := services.NewAuthService(userStore, sessionStore)
	authHandler := handlers.NewAuthHandler(authService)

	api.Post("/auth/signup", authHandler.Signup)
	api.Post("/auth/login", authHandler.Login)

	protected := api.Group("", middleware.AuthRequired(sessionStore))

	protected.Post("/auth/logout", authHandler.Logout)
	protected.Get("/auth/me", authHandler.GetCurrentUser)

	// Projects
	projectStore := store.NewProjectStore()
	projectService := services.NewProjectService(projectStore)
	projectHandler := handlers.NewProjectHandler(projectService)

	protected.Get("/projects", projectHandler.ListProjects)
	protected.Get("/projects/:id", projectHandler.GetProject)
	protected.Post("/projects", projectHandler.CreateProject)
	protected.Patch("/projects/:id", projectHandler.UpdateProject)
	protected.Delete("/projects/:id", projectHandler.DeleteProject)

	// Project members
	projectMemberStore := store.NewProjectMemberStore()
	projectMemberService := services.NewProjectMemberService(projectMemberStore)
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberService)

	protected.Get("/projects/:projectId/members", projectMemberHandler.ListProjectMembers)
	protected.Post("/projects/:projectId/members", projectMemberHandler.AddProjectMember)
	protected.Patch("/projects/:projectId/members/:userId", projectMemberHandler.UpdateProjectMemberRole)
	protected.Delete("/projects/:projectId/members/:userId", projectMemberHandler.RemoveProjectMember)

	// Priorities
	priorityStore := store.NewPriorityStore()
	priorityService := services.NewPriorityService(priorityStore, projectStore)
	priorityHandler := handlers.NewPriorityHandler(priorityService)

	protected.Get("/projects/:projectID/priorities", priorityHandler.ListPriorities)
	protected.Post("/projects/:projectID/priorities", priorityHandler.CreatePriority)
	protected.Patch("/priorities/:priorityID", priorityHandler.UpdatePriority)
	protected.Patch("/priorities/:priorityID/position", priorityHandler.UpdatePriorityPosition)
	protected.Delete("/priorities/:priorityID", priorityHandler.DeletePriority)

	// Statuses
	statusStore := store.NewStatusStore()
	statusService := services.NewStatusService(statusStore, projectStore)
	statusHandler := handlers.NewStatusHandler(statusService)

	protected.Get("/projects/:projectID/statuses", statusHandler.ListStatuses)
	protected.Post("/projects/:projectID/statuses", statusHandler.CreateStatus)
	protected.Patch("/statuses/:statusID", statusHandler.UpdateStatus)
	protected.Patch("/statuses/:statusID/position", statusHandler.UpdateStatusPosition)
	protected.Delete("/statuses/:statusID", statusHandler.DeleteStatus)

	// Tags
	tagStore := store.NewTagStore()
	tagService := services.NewTagService(tagStore, projectStore)
	tagHandler := handlers.NewTagHandler(tagService)

	protected.Get("/projects/:projectID/tags", tagHandler.ListTags)
	protected.Post("/projects/:projectID/tags", tagHandler.CreateTag)
	protected.Patch("/tags/:tagID", tagHandler.UpdateTag)
	protected.Delete("/tags/:tagID", tagHandler.DeleteTag)

	// Tickets
	ticketStore := store.NewTicketStore()
	ticketService := services.NewTicketService(
		ticketStore,
		projectStore,
		projectMemberStore,
		statusStore,
		priorityStore,
		userStore,
		tagStore,
	)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	protected.Get("/projects/:projectID/tickets", ticketHandler.ListProjectTickets)
	protected.Post("/projects/:projectID/tickets", ticketHandler.CreateProjectTicket)
	protected.Get("/tickets/:ticketID", ticketHandler.GetTicket)
	protected.Patch("/tickets/:ticketID", ticketHandler.UpdateTicket)
	protected.Delete("/tickets/:ticketID", ticketHandler.DeleteTicket)

	// Comments
	commentStore := store.NewCommentStore()
	commentService := services.NewCommentService(commentStore, projectMemberStore)
	commentHandler := handlers.NewCommentHandler(commentService)

	protected.Get("/tickets/:ticketID/comments", commentHandler.ListComments)
	protected.Post("/tickets/:ticketID/comments", commentHandler.CreateComment)
	protected.Patch("/comments/:commentID", commentHandler.UpdateComment)
	protected.Delete("/comments/:commentID", commentHandler.DeleteComment)

	// Activity
	ticketActivityStore := store.NewTicketActivityStore()
	ticketActivityService := services.NewTicketActivityService(ticketActivityStore, projectMemberStore)
	ticketActivityHandler := handlers.NewTicketActivityHandler(ticketActivityService)

	protected.Get("/tickets/:ticketID/activity", ticketActivityHandler.GetTicketActivity)
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
