package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ListPriorities(c *fiber.Ctx) error {
	// Check if projectID in path
	projectIDStr := c.Params("id")

	var projectID *int
	if projectIDStr != "" {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		projectID = &id
	}

	priorities, err := models.ListPriorities(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(priorities)
}

func CreatePriority(c *fiber.Ctx) error {
	// Check if projectID in path
	projectIDStr := c.Params("id")

	var projectID *int
	if projectIDStr != "" {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		projectID = &id
	}

	var req struct {
		Label string `json:"label"`
		Color string `json:"color"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Label == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Label is required",
		})
	}

	priority, err := models.CreatePriority(projectID, req.Label, req.Color)
	if err != nil {
		if err.Error() == "invalid color format: must be #RRGGBB" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid color format",
			})
		}
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Project not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(priority)
}

func UpdatePriority(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid priority ID",
		})
	}

	var req struct {
		Label *string `json:"label"`
		Color *string `json:"color"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Label == nil && req.Color == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	priority, err := models.UpdatePriority(id, req.Label, req.Color)
	if err != nil {
		if errors.Is(err, models.ErrPriorityNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Priority not found",
			})
		}
		if err.Error() == "invalid color format: must be #RRGGBB" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid color format",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(priority)
}

func UpdatePriorityPosition(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid priority ID",
		})
	}

	var req struct {
		Position int `json:"position"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	priority, err := models.UpdatePriorityPosition(id, req.Position)
	if err != nil {
		if errors.Is(err, models.ErrPriorityNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Priority not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(priority)
}

func DeletePriority(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid priority ID",
		})
	}

	if err := models.DeletePriority(id); err != nil {
		if errors.Is(err, models.ErrPriorityNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Priority not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}
