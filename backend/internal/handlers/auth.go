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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, err := models.CreateUser(req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, models.ErrEmailExists) {
			return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
		}
		if errors.Is(err, models.ErrInvalidEmail) ||
			errors.Is(err, models.ErrInvalidPassword) ||
			errors.Is(err, models.ErrInvalidName) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(500)
	}

	if err := createSessionAndSetCookie(c, user.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create session"})
	}

	return c.Status(201).JSON(user)
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, err := models.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := createSessionAndSetCookie(c, user.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create session"})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		_ = models.RevokeSession(sessionID)
	}

	c.ClearCookie("session_id")
	return c.SendStatus(204)
}

func GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(401)
	}
	return c.JSON(user)
}

func createSessionAndSetCookie(c *fiber.Ctx, userID int) error {
	session, err := models.CreateSession(userID, c.IP(), string(c.Request().Header.UserAgent()))
	if err != nil {
		return err
	}

	// Local development is HTTP, so Secure cookies must be disabled there.
	isDevelopment := strings.EqualFold(os.Getenv("ENV"), "development")
	useSecureCookie := !isDevelopment

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		HTTPOnly: true,
		Secure:   useSecureCookie,
		SameSite: "Lax",
		Path:     "/",
	})

	return nil
}
