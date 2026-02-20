package handlers

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ListComments(c *fiber.Ctx) error {
	ticketID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	comments, err := models.ListComments(ticketID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(comments)
}

func CreateComment(c *fiber.Ctx) error {
	ticketID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	var req struct {
		AuthorID int    `json:"author_id"`
		Body     string `json:"body"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Comment body is required",
		})
	}
	if req.AuthorID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Author ID is required",
		})
	}

	comment, err := models.CreateComment(ticketID, req.AuthorID, req.Body)
	if err != nil {
		if errors.Is(err, models.ErrTicketNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Ticket not found",
			})
		}
		if errors.Is(err, models.ErrUserNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Author not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(comment)
}

func UpdateComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	var req struct {
		Body string `json:"body"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Comment body is required",
		})
	}

	comment, err := models.UpdateComment(id, req.Body)
	if err != nil {
		if errors.Is(err, models.ErrCommentNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(comment)
}

func DeleteComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	if err := models.DeleteComment(id); err != nil {
		if errors.Is(err, models.ErrCommentNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}
