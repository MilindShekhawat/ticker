package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

type AddMemberRequest struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
}

type UpdateMemberRoleRequest struct {
	Role int `json:"role"`
}

func ListProjectMembers(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectId")
	if err != nil {
		return c.SendStatus(400)
	}

	members, err := models.ListProjectMembers(projectID)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(members)
}

func AddProjectMember(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectId")
	if err != nil {
		return c.SendStatus(400)
	}

	var req AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(400)
	}

	err = models.AddProjectMember(projectID, req.UserID, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidRole):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, models.ErrMemberAlreadyExists):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, models.ErrOwnerAlreadyExists):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.SendStatus(500)
		}
	}

	return c.SendStatus(201)
}

func UpdateProjectMemberRole(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectId")
	if err != nil {
		return c.SendStatus(400)
	}

	userID, err := c.ParamsInt("userId")
	if err != nil {
		return c.SendStatus(400)
	}

	var req UpdateMemberRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(400)
	}

	err = models.UpdateProjectMemberRole(projectID, userID, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidRole),
			errors.Is(err, models.ErrCannotDemoteOwner),
			errors.Is(err, models.ErrCannotPromoteToOwner):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, models.ErrNotProjectMember):
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.SendStatus(500)
		}
	}

	return c.SendStatus(200)
}

func RemoveProjectMember(c *fiber.Ctx) error {
	projectID, err := c.ParamsInt("projectId")
	if err != nil {
		return c.SendStatus(400)
	}

	userID, err := c.ParamsInt("userId")
	if err != nil {
		return c.SendStatus(400)
	}

	err = models.RemoveProjectMember(projectID, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotProjectMember):
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, models.ErrCannotRemoveOwner):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.SendStatus(500)
		}
	}

	return c.SendStatus(204)
}
