package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type PriorityHandler struct {
	service services.PriorityService
}

func NewPriorityHandler(service services.PriorityService) *PriorityHandler {
	return &PriorityHandler{service: service}
}

type CreatePriorityRequest struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

type UpdatePriorityRequest struct {
	Label *string `json:"label"`
	Color *string `json:"color"`
}

type UpdatePriorityPositionRequest struct {
	Position int `json:"position"`
}

func (h *PriorityHandler) ListPriorities(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	priorities, err := h.service.List(projectID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(priorities)
}

func (h *PriorityHandler) CreatePriority(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req CreatePriorityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	priority, err := h.service.Create(projectID, req.Label, req.Color)
	if err != nil {
		return priorityErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(priority)
}

func (h *PriorityHandler) UpdatePriority(c *fiber.Ctx) error {
	priorityID, err := c.ParamsInt("priorityID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid priority ID"})
	}

	var req UpdatePriorityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Label == nil && req.Color == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	priority, err := h.service.Update(priorityID, req.Label, req.Color)
	if err != nil {
		return priorityErrorResponse(c, err)
	}

	return c.JSON(priority)
}

func (h *PriorityHandler) UpdatePriorityPosition(c *fiber.Ctx) error {
	priorityID, err := c.ParamsInt("priorityID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid priority ID"})
	}

	var req UpdatePriorityPositionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	priority, err := h.service.UpdatePosition(priorityID, req.Position)
	if err != nil {
		return priorityErrorResponse(c, err)
	}

	return c.JSON(priority)
}

func (h *PriorityHandler) DeletePriority(c *fiber.Ctx) error {
	priorityID, err := c.ParamsInt("priorityID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid priority ID"})
	}

	if err := h.service.Delete(priorityID); err != nil {
		return priorityErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func priorityErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, store.ErrPriorityNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Priority not found"})
	case errors.Is(err, store.ErrProjectNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	case errors.Is(err, services.ErrLabelRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Label is required"})
	case errors.Is(err, services.ErrInvalidColor):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid color format: must be #RRGGBB"})
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
