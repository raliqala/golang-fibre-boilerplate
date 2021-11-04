package welcome

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/controllers"
)

func Welcome(router fiber.Router) {
	router.Get("/", controllers.Welcome)
}
