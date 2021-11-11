package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/controllers/users"
)

func ProfileRoutes(router fiber.Router) {
	router.Post("/profile", users.GetProfile)
}
