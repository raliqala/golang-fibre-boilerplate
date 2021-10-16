package welcomeRoute

import (
	"github.com/gofiber/fiber/v2"
	WelcomeController "github.com/raliqala/safepass_api/src/controllers"
)

func Welcome(router fiber.Router) {
	router.Get("/", WelcomeController.Welcome)
}
