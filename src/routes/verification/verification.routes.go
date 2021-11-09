package verification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/controllers"
)

func VerifyRoute(router fiber.Router) {
	router.Post("/verify/:token", controllers.EmailVerification)
}
