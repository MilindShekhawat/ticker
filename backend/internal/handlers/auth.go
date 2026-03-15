package handlers

import (
	"errors"
	"os"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, session, err := h.service.Signup(
		req.Email,
		req.Password,
		req.Name,
		c.IP(),
		string(c.Request().Header.UserAgent()),
	)
	if err != nil {
		return authErrorResponse(c, err)
	}

	setSessionCookie(c, session.ID)
	return c.Status(201).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, session, err := h.service.Login(
		req.Email,
		req.Password,
		c.IP(),
		string(c.Request().Header.UserAgent()),
	)
	if err != nil {
		return authErrorResponse(c, err)
	}

	setSessionCookie(c, session.ID)
	return c.JSON(user)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		err := h.service.Logout(sessionID)
		if err != nil && !errors.Is(err, store.ErrInvalidSession) {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	c.ClearCookie("session_id")
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.JSON(user)
}

func setSessionCookie(c *fiber.Ctx, sessionID string) {
	// Local development is HTTP, so Secure cookies must be disabled there.
	isDevelopment := strings.EqualFold(os.Getenv("ENV"), "development")
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HTTPOnly: true,
		Secure:   !isDevelopment,
		SameSite: "Lax",
		Path:     "/",
	})
}

func authErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, services.ErrInvalidEmail),
		errors.Is(err, services.ErrInvalidPassword),
		errors.Is(err, services.ErrInvalidName):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, store.ErrEmailExists):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already in use"})
	case errors.Is(err, store.ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to process request"})
	}
}
