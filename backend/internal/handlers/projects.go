package handlers

import (
	"errors"
	"strings"
	"unicode"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func isValidKeyPrefix(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Projects

func ListProjects(c *fiber.Ctx) error {
	projects, err := models.ListProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(projects)
}

func CreateProject(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		KeyPrefix   string `json:"key_prefix"`
		CreatedBy   int    `json:"created_by"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name is required",
		})
	}
	if len(req.Name) > 100 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name too long (max 100 characters)",
		})
	}
	if req.KeyPrefix == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Key prefix is required",
		})
	}
	if len(req.KeyPrefix) > 5 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Key prefix too long (max 5 characters)",
		})
	}
	if !isValidKeyPrefix(req.KeyPrefix) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Key prefix must be uppercase letters and numbers only",
		})
	}
	if req.CreatedBy == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Created by is required",
		})
	}

	project, err := models.CreateProject(req.Name, req.Description, strings.ToUpper(req.KeyPrefix), req.CreatedBy)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateKeyPrefix) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Key prefix already exists",
			})
		}
		if errors.Is(err, models.ErrUserNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid created by user ID",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(project)
}

func GetProject(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	project, err := models.GetProjectByID(id)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Project not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(project)
}

func UpdateProject(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == nil && req.Description == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	if req.Name != nil {
		if *req.Name == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Name cannot be empty",
			})
		}
		if len(*req.Name) > 100 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Name too long (max 100 characters)",
			})
		}
	}

	project, err := models.UpdateProject(id, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Project not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(project)
}

func DeleteProject(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	if err := models.DeleteProject(id); err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Project not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}
