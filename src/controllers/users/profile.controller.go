package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/database"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
)

func GetProfile(c *fiber.Ctx) error {
	db := database.DB
	id := c.Locals("id")

	u := new(models.User)
	if res := db.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Cannot find the User"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    u,
	})
}
