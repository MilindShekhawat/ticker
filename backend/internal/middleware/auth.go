package middleware

import (
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired(sessionStore store.SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		user, err := sessionStore.Validate(sessionID)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user", user)
		return c.Next()
	}
}
