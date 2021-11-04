package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/controllers/users"
)

func AuthRoutes(router fiber.Router) {
	user := router.Group("/u")
	user.Post("/signup", users.SignUp)
	user.Post("/signin", users.SignIn)
	user.Get("/access-token", users.GetAccessToken)
}

// func ProfileRoutes(router fiber.Router) {
// 	router.Get("/profile", users.GetProfile)
// }
