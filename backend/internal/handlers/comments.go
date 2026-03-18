package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/services"
	"github.com/MilindShekhawat/ticker/internal/store"
	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	service services.CommentService
}

func NewCommentHandler(service services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

type CreateCommentRequest struct {
	Body string `json:"body"`
}

type UpdateCommentRequest struct {
	Body string `json:"body"`
}

func (h *CommentHandler) ListComments(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	ticketID, err := c.ParamsInt("ticketID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	comments, err := h.service.ListForUser(ticketID, user.ID)
	if err != nil {
		return commentErrorResponse(c, err)
	}

	return c.JSON(comments)
}

func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	ticketID, err := c.ParamsInt("ticketID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	comment, err := h.service.CreateForUser(ticketID, user.ID, req.Body)
	if err != nil {
		return commentErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	commentID, err := c.ParamsInt("commentID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	var req UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	comment, err := h.service.UpdateForUser(commentID, user.ID, req.Body)
	if err != nil {
		return commentErrorResponse(c, err)
	}

	return c.JSON(comment)
}

func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	commentID, err := c.ParamsInt("commentID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	if err := h.service.DeleteForUser(commentID, user.ID); err != nil {
		return commentErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func commentErrorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, services.ErrInvalidCommentBody):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comment body is required"})
	case errors.Is(err, store.ErrTicketNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ticket not found"})
	case errors.Is(err, store.ErrCommentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	case errors.Is(err, store.ErrNotProjectMember):
		return c.SendStatus(fiber.StatusForbidden)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
