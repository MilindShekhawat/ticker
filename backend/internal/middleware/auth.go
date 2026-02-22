package middleware

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

const contextUserKey = "user"

func AuthRequired(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.SendStatus(401)
	}

	user, err := models.ValidateSession(sessionID)
	if err != nil {
		if errors.Is(err, models.ErrInvalidSession) {
			c.ClearCookie("session_id")
			return c.SendStatus(401)
		}
		return c.SendStatus(500)
	}

	c.Locals(contextUserKey, user)
	return c.Next()
}
