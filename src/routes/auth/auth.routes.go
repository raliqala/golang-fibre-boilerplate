package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/controllers/users"
)

func AuthRoutes(router fiber.Router) {
	user := router.Group("/u")
	user.Post("/", users.SignUp)
}
