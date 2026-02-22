package handlers

import (
	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetTicketActivity(c *fiber.Ctx) error {
	ticketID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ticket ID",
		})
	}

	activities, err := models.ListActivityByTicket(ticketID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(activities)
}
