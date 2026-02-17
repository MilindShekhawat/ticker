package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetTickets(c *fiber.Ctx) error {
	projectId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	tickets, err := models.GetTicketsByProjectID(projectId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tickets)
}

func GetTicket(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	ticket, err := models.GetTicketByID(id)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Ticket not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ticket)
}

func CreateTicket(c *fiber.Ctx) error {
	projectId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StatusID    int    `json:"status_id"`
		PriorityID  int    `json:"priority_id"`
		CreatedBy   int    `json:"created_by"`
		AssigneeID  *int   `json:"assignee_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Title is required",
		})
	}
	if len(req.Title) > 200 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Title too long (max 200 characters)",
		})
	}
	if req.StatusID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Status ID is required",
		})
	}
	if req.PriorityID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Priority ID is required",
		})
	}
	if req.CreatedBy == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Created by is required",
		})
	}

	ticket, err := models.CreateTicket(
		projectId,
		req.Title,
		req.Description,
		req.StatusID,
		req.PriorityID,
		req.CreatedBy,
		req.AssigneeID,
	)

	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		if errors.Is(err, models.ErrStatusNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid status ID",
			})
		}
		if errors.Is(err, models.ErrPriorityNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid priority ID",
			})
		}
		if errors.Is(err, models.ErrUserNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid created by user ID",
			})
		}
		if errors.Is(err, models.ErrAssigneeNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid assignee ID",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(ticket)
}

func UpdateTicket(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		StatusID    *int    `json:"status_id"`
		PriorityID  *int    `json:"priority_id"`
		AssigneeID  *int    `json:"assignee_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == nil && req.Description == nil && req.StatusID == nil && req.PriorityID == nil && req.AssigneeID == nil {
		return c.Status(400).JSON(fiber.Map{"error": "No fields to update"})
	}
	if req.Title != nil && *req.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title cannot be empty"})
	}
	if req.Title != nil && len(*req.Title) > 200 {
		return c.Status(400).JSON(fiber.Map{"error": "Title too long"})
	}
	if req.StatusID != nil && *req.StatusID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Status ID cannot be 0"})
	}
	if req.PriorityID != nil && *req.PriorityID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Priority ID cannot be 0"})
	}

	ticket, err := models.UpdateTicket(id, req.Title, req.Description, req.StatusID, req.PriorityID, req.AssigneeID)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Ticket not found",
			})
		}
		if errors.Is(err, models.ErrStatusNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid status ID",
			})
		}
		if errors.Is(err, models.ErrPriorityNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid priority ID",
			})
		}
		if errors.Is(err, models.ErrAssigneeNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid assignee ID",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ticket)
}

func DeleteTicket(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	if err := models.DeleteTicket(id); err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Ticket not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}
