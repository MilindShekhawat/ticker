package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	service services.TagService
}

func NewTagHandler(service services.TagService) *TagHandler {
	return &TagHandler{service: service}
}

type CreateTagRequest struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

type UpdateTagRequest struct {
	Label *string `json:"label"`
	Color *string `json:"color"`
}

func (h *TagHandler) ListTags(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	tags, err := h.service.List(projectID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(tags)
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tag, err := h.service.Create(projectID, req.Label, req.Color)
	if err != nil {
		return tagErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(tag)
}

func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	tagID, err := c.ParamsInt("tagID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag ID"})
	}

	var req UpdateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Label == nil && req.Color == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	tag, err := h.service.Update(tagID, req.Label, req.Color)
	if err != nil {
		return tagErrorResponse(c, err)
	}

	return c.JSON(tag)
}

func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	tagID, err := c.ParamsInt("tagID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag ID"})
	}

	if err := h.service.Delete(tagID); err != nil {
		return tagErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func tagErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, store.ErrTagNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tag not found"})
	case errors.Is(err, store.ErrProjectNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	case errors.Is(err, services.ErrInvalidColor):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid color format: must be #RRGGBB"})
	case errors.Is(err, services.ErrLabelRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Label is required"})
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
