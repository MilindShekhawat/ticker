package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	KeyPrefix   string `json:"key_prefix"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func ListProjects(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projects, err := models.ListProjectsByUser(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(projects)
}

func CreateProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req CreateProjectRequest

	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	project, err := models.CreateProject(
		req.Name,
		req.Description,
		req.KeyPrefix,
		user.ID,
	)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidProject),
			errors.Is(err, models.ErrInvalidProjectName),
			errors.Is(err, models.ErrInvalidKeyPrefix),
			errors.Is(err, models.ErrDuplicateKeyPrefix),
			errors.Is(err, models.ErrUserNotFound):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func GetProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	project, err := models.GetProjectByUser(id, user.ID)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(project)
}

func UpdateProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var req UpdateProjectRequest

	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	project, err := models.UpdateProjectByUser(id, user.ID, req.Name, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidProjectName),
			errors.Is(err, models.ErrInvalidProjectUpdate):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		case errors.Is(err, models.ErrProjectNotFound):
			return c.SendStatus(fiber.StatusNotFound)

		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.JSON(project)
}

func DeleteProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := models.DeleteProjectByUser(id, user.ID); err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
