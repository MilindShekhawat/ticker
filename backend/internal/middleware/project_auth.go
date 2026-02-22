package middleware

import (
	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func RequireProjectMember() fiber.Handler {
	return func(c *fiber.Ctx) error {
		projectID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
		}

		userID := c.Locals("user_id").(int)

		err = models.IsProjectMember(projectID, userID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Not a project member"})
		}

		return c.Next()
	}
}

func RequireProjectRole(minRole int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		projectID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
		}

		userID := c.Locals("user_id").(int)

		role, err := models.GetUserRole(projectID, userID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Not authorized"})
		}

		if role > minRole {
			return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
		}

		return c.Next()
	}
}
