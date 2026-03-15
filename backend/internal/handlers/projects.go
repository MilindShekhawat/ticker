package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	service services.ProjectService
}

func NewProjectHandler(service services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	KeyPrefix   string `json:"key_prefix"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projects, err := h.service.ListForUser(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(projects)
}

func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	project, err := h.service.GetForUser(id, user.ID)
	if err != nil {
		return projectErrorResponse(c, err)
	}

	return c.JSON(project)
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	project, err := h.service.CreateForUser(
		req.Name,
		req.Description,
		req.KeyPrefix,
		user.ID,
	)
	if err != nil {
		return projectErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
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

	project, err := h.service.UpdateForUser(id, user.ID, req.Name, req.Description)
	if err != nil {
		return projectErrorResponse(c, err)
	}

	return c.JSON(project)
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.DeleteForUser(id, user.ID); err != nil {
		return projectErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func projectErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, store.ErrProjectNotFound):
		return c.SendStatus(fiber.StatusNotFound)
	case errors.Is(err, services.ErrInvalidProject),
		errors.Is(err, services.ErrInvalidProjectName),
		errors.Is(err, services.ErrInvalidKeyPrefix),
		errors.Is(err, services.ErrInvalidProjectUpdate),
		errors.Is(err, store.ErrDuplicateKeyPrefix):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
