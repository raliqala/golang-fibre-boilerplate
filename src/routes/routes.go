package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/raliqala/safepass_api/src/routes/auth"
	"github.com/raliqala/safepass_api/src/routes/welcome"
)

func RouteSetup(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// protected := app.Group("/private")

	// protected.Use(middleware.SecureAuth)
	welcome.Welcome(api)
	auth.AuthRoutes(api)
	// auth.ProfileRoutes(protected)
}
