package routes

import (
	"github.com/gofiber/fiber/v2"
	WelcomeRoutes "github.com/raliqala/safepass_api/src/routes/welcome"
)

func RouteSetup(app *fiber.App) {
	api := app.Group("/api")
	WelcomeRoutes.Welcome(api)
}
