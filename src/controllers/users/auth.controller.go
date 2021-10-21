package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/database"
	"github.com/raliqala/safepass_api/src/models"
)

func SignUp(c *fiber.Ctx) error {
	// testing
	db := database.DB
	data := new(models.Author)

	err := c.BodyParser(data)

	// log.Fatalln(data, err)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	err = db.Create(&data).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create author", "data": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"test":    data,
	})
}
