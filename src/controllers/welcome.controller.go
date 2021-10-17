package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/services"
)

func Welcome(c *fiber.Ctx) error {
	accessToken, refreshToken := services.GenerateTokens("14d1ea63-fb27-4eb8-964c-c49da50939cb")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"message":       "All good here!.. 🚀 ",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
