package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ListTags(c *fiber.Ctx) error {
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

	tags, err := models.ListTags(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tags)
}

func CreateTag(c *fiber.Ctx) error {
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

	tag, err := models.CreateTag(projectID, req.Label, req.Color)
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

	return c.Status(201).JSON(tag)
}

func UpdateTag(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid tag ID",
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

	tag, err := models.UpdateTag(id, req.Label, req.Color)
	if err != nil {
		if errors.Is(err, models.ErrTagNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Tag not found",
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

	return c.JSON(tag)
}

func DeleteTag(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid tag ID",
		})
	}

	if err := models.DeleteTag(id); err != nil {
		if errors.Is(err, models.ErrTagNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Tag not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}
