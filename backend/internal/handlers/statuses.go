package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ListStatuses(c *fiber.Ctx) error {
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

	statuses, err := models.ListStatuses(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(statuses)
}

func CreateStatus(c *fiber.Ctx) error {
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

	status, err := models.CreateStatus(projectID, req.Label, req.Color)
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

	return c.Status(201).JSON(status)
}

func UpdateStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid status ID",
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

	status, err := models.UpdateStatus(id, req.Label, req.Color)
	if err != nil {
		if errors.Is(err, models.ErrStatusNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Status not found",
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

	return c.JSON(status)
}

func UpdateStatusPosition(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid status ID",
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

	status, err := models.UpdateStatusPosition(id, req.Position)
	if err != nil {
		if errors.Is(err, models.ErrStatusNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Status not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(status)
}

func DeleteStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid status ID",
		})
	}

	if err := models.DeleteStatus(id); err != nil {
		if errors.Is(err, models.ErrStatusNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Status not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}
