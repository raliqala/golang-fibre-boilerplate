package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/raliqala/golang-fibre-boilerplate/src/routes/auth"
	"github.com/raliqala/golang-fibre-boilerplate/src/routes/verification"
	"github.com/raliqala/golang-fibre-boilerplate/src/routes/welcome"
)

func RouteSetup(app *fiber.App) {
	api := app.Group("/api", logger.New())
	v1 := api.Group("auth")

	// protected := app.Group("/private")

	// protected.Use(middleware.SecureAuth)
	welcome.Welcome(v1)
	auth.AuthRoutes(v1)
	verification.VerifyRoute(v1)
	// auth.ProfileRoutes(protected)
}
