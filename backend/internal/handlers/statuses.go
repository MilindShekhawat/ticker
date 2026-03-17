package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type StatusHandler struct {
	service services.StatusService
}

func NewStatusHandler(service services.StatusService) *StatusHandler {
	return &StatusHandler{service: service}
}

type CreateStatusRequest struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

type UpdateStatusRequest struct {
	Label *string `json:"label"`
	Color *string `json:"color"`
}

type UpdateStatusPositionRequest struct {
	Position int `json:"position"`
}

func (h *StatusHandler) ListStatuses(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	statuses, err := h.service.List(projectID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(statuses)
}

func (h *StatusHandler) CreateStatus(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req CreateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	status, err := h.service.Create(projectID, req.Label, req.Color)
	if err != nil {
		return statusErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(status)
}

func (h *StatusHandler) UpdateStatus(c *fiber.Ctx) error {
	statusID, err := c.ParamsInt("statusID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status ID"})
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Label == nil && req.Color == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	status, err := h.service.Update(statusID, req.Label, req.Color)
	if err != nil {
		return statusErrorResponse(c, err)
	}

	return c.JSON(status)
}

func (h *StatusHandler) UpdateStatusPosition(c *fiber.Ctx) error {
	statusID, err := c.ParamsInt("statusID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status ID"})
	}

	var req UpdateStatusPositionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	status, err := h.service.UpdatePosition(statusID, req.Position)
	if err != nil {
		return statusErrorResponse(c, err)
	}

	return c.JSON(status)
}

func (h *StatusHandler) DeleteStatus(c *fiber.Ctx) error {
	statusID, err := c.ParamsInt("statusID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status ID"})
	}

	if err := h.service.Delete(statusID); err != nil {
		return statusErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func statusErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, store.ErrStatusNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Status not found"})
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
