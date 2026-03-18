package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	service services.TicketService
}

func NewTicketHandler(service services.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusID    int    `json:"status_id"`
	PriorityID  int    `json:"priority_id"`
	AssigneeID  *int   `json:"assignee_id"`
	TagIDs      []int  `json:"tag_ids,omitempty"`
}

type UpdateTicketRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	StatusID    *int    `json:"status_id"`
	PriorityID  *int    `json:"priority_id"`
	AssigneeID  *int    `json:"assignee_id"`
	TagIDs      *[]int  `json:"tag_ids"`
}

func (h *TicketHandler) GetTicket(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	ticketID, err := c.ParamsInt("ticketID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	ticket, err := h.service.GetForUser(ticketID, user.ID)
	if err != nil {
		return ticketErrorResponse(c, err)
	}

	return c.JSON(ticket)
}

func (h *TicketHandler) ListProjectTickets(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	tickets, err := h.service.ListForUser(projectID, user.ID)
	if err != nil {
		return ticketErrorResponse(c, err)
	}

	return c.JSON(tickets)
}

func (h *TicketHandler) CreateProjectTicket(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ticket, err := h.service.CreateForUser(
		projectID,
		req.Title,
		req.Description,
		req.StatusID,
		req.PriorityID,
		user.ID,
		req.AssigneeID,
		req.TagIDs,
	)
	if err != nil {
		return ticketErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(ticket)
}

func (h *TicketHandler) UpdateTicket(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	ticketID, err := c.ParamsInt("ticketID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	var req UpdateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Title == nil && req.Description == nil && req.StatusID == nil &&
		req.PriorityID == nil && req.AssigneeID == nil && req.TagIDs == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	ticket, err := h.service.UpdateForUser(
		ticketID,
		user.ID,
		req.Title,
		req.Description,
		req.StatusID,
		req.PriorityID,
		req.AssigneeID,
		req.TagIDs,
	)
	if err != nil {
		return ticketErrorResponse(c, err)
	}

	return c.JSON(ticket)
}

func (h *TicketHandler) DeleteTicket(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	ticketID, err := c.ParamsInt("ticketID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	if err := h.service.DeleteForUser(ticketID, user.ID); err != nil {
		return ticketErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func ticketErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, services.ErrTicketTitleRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	case errors.Is(err, services.ErrTicketTitleTooLong):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title too long (max 200 characters)"})
	case errors.Is(err, services.ErrInvalidStatusID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status ID"})
	case errors.Is(err, services.ErrInvalidPriorityID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid priority ID"})
	case errors.Is(err, services.ErrInvalidAssigneeID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignee ID"})
	case errors.Is(err, store.ErrProjectNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	case errors.Is(err, store.ErrTicketNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	case errors.Is(err, store.ErrStatusNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status ID"})
	case errors.Is(err, store.ErrPriorityNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid priority ID"})
	case errors.Is(err, store.ErrTagNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag ID"})
	case errors.Is(err, store.ErrUserNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid assignee ID"})
	case errors.Is(err, store.ErrNotProjectMember):
		return c.SendStatus(fiber.StatusForbidden)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
