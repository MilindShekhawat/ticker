package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type TicketActivityHandler struct {
	service services.TicketActivityService
}

func NewTicketActivityHandler(service services.TicketActivityService) *TicketActivityHandler {
	return &TicketActivityHandler{service: service}
}

func (h *TicketActivityHandler) GetTicketActivity(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	ticketID, err := c.ParamsInt("ticketID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	activities, err := h.service.ListForUser(ticketID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTicketNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
		case errors.Is(err, store.ErrNotProjectMember):
			return c.SendStatus(fiber.StatusForbidden)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.JSON(activities)
}
