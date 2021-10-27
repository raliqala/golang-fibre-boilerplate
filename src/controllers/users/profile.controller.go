package users

import (
	"github.com/gofiber/fiber"
	"github.com/raliqala/safepass_api/src/database"
	"github.com/raliqala/safepass_api/src/models"
)

func GetUserData(c *fiber.Ctx) error {
	db := database.DB
	id := c.Locals("id")

	u := new(models.User)
	if res := db.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}
