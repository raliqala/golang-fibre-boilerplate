package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/controllers/users"
)

func AuthRoutes(router fiber.Router) {
	router.Post("/signup", users.SignUp)
	router.Post("/signin", users.SignIn)
	router.Get("/access-token", users.GetAccessToken)
}
