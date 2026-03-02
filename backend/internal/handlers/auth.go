package handlers

import (
	"errors"
	"os"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *fiber.Ctx) error {
	var req SignupRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, session, err := models.SignupWithSession(
		req.Email,
		req.Password,
		req.Name,
		c.IP(),
		string(c.Request().Header.UserAgent()),
	)
	if err != nil {
		return handleSignupError(c, err)
	}

	setSessionCookie(c, session.ID)

	return c.Status(201).JSON(user)
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, session, err := models.LoginWithSession(
		req.Email,
		req.Password,
		c.IP(),
		string(c.Request().Header.UserAgent()),
	)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to process request"})
	}

	setSessionCookie(c, session.ID)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")

	if sessionID != "" {
		err := models.RevokeSession(sessionID)
		if err != nil && !errors.Is(err, models.ErrInvalidSession) {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	c.ClearCookie("session_id")
	return c.SendStatus(fiber.StatusNoContent)
}

func GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.JSON(user)
}

func setSessionCookie(c *fiber.Ctx, sessionID string) {
	// Local development is HTTP, so Secure cookies must be disabled there.
	isDevelopment := strings.EqualFold(os.Getenv("ENV"), "development")
	useSecureCookie := !isDevelopment

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HTTPOnly: true,
		Secure:   useSecureCookie,
		SameSite: "Lax",
		Path:     "/",
	})
}

func handleSignupError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrEmailExists):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already in use"})
	case errors.Is(err, models.ErrInvalidEmail),
		errors.Is(err, models.ErrInvalidPassword),
		errors.Is(err, models.ErrInvalidName):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid signup details"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to process request"})
	}
}
