package middleware

import (
	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user, err := models.ValidateSession(sessionID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}

	c.Locals("user", user)
	c.Locals("user_id", user.ID)

	return c.Next()
}
