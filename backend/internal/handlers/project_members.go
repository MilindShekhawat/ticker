package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type ProjectMemberHandler struct {
	service services.ProjectMemberService
}

func NewProjectMemberHandler(service services.ProjectMemberService) *ProjectMemberHandler {
	return &ProjectMemberHandler{service: service}
}

type AddMemberRequest struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
}

type UpdateMemberRoleRequest struct {
	Role int `json:"role"`
}

func (h *ProjectMemberHandler) ListProjectMembers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID, err := projectIDFromParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	members, err := h.service.ListForUser(projectID, user.ID)
	if err != nil {
		return projectMemberErrorResponse(c, err)
	}

	return c.JSON(members)
}

func (h *ProjectMemberHandler) AddProjectMember(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID, err := projectIDFromParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.AddForUser(projectID, user.ID, req.UserID, req.Role); err != nil {
		return projectMemberErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *ProjectMemberHandler) UpdateProjectMemberRole(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID, err := projectIDFromParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	memberUserID, err := userIDFromParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req UpdateMemberRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.UpdateRoleForUser(projectID, user.ID, memberUserID, req.Role); err != nil {
		return projectMemberErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *ProjectMemberHandler) RemoveProjectMember(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	projectID, err := projectIDFromParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	memberUserID, err := userIDFromParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.service.RemoveForUser(projectID, user.ID, memberUserID); err != nil {
		return projectMemberErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func projectMemberErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, services.ErrInvalidRole):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
	case errors.Is(err, store.ErrMemberAlreadyExists):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Member already exists"})
	case errors.Is(err, store.ErrOwnerAlreadyExists):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Owner already exists"})
	case errors.Is(err, store.ErrCannotDemoteOwner):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot demote project owner"})
	case errors.Is(err, store.ErrCannotPromoteToOwner):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Use transfer ownership instead"})
	case errors.Is(err, store.ErrCannotRemoveOwner):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot remove project owner"})
	case errors.Is(err, services.ErrInsufficientProjectRole):
		return c.SendStatus(fiber.StatusForbidden)
	case errors.Is(err, store.ErrNotProjectMember):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project member not found"})
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func projectIDFromParams(c *fiber.Ctx) (int, error) {
	projectID, err := c.ParamsInt("projectID")
	if err == nil {
		return projectID, nil
	}
	return c.ParamsInt("projectId")
}

func userIDFromParams(c *fiber.Ctx) (int, error) {
	userID, err := c.ParamsInt("userID")
	if err == nil {
		return userID, nil
	}
	return c.ParamsInt("userId")
}
