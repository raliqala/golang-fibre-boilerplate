package welcome

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/controllers"
)

func Welcome(router fiber.Router) {
	router.Get("/", controllers.Welcome)
}
