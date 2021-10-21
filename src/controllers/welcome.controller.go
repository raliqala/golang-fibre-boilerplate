package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/services"
)

func Welcome(c *fiber.Ctx) error {
	services.GenerateTokens("d7fa0111-7371-41a2-8a3f-55a8c509dd61")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "All good here!.. 🚀 ",
	})
}
