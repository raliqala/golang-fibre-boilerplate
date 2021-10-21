package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/config"
	"github.com/raliqala/safepass_api/src/database"
	"github.com/raliqala/safepass_api/src/routes"
)

func main() {
	fibreConfig := fiber.Config{
		ServerHeader: config.Config("APP_NAME"),
	}
	app := fiber.New(fibreConfig)

	database.ConnectPg()

	// declaration
	// csrfConfig := csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Strict",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }

	// app config
	// app.Use(logger.New())
	// app.Use(cors.New())
	// app.Use(csrf.New(csrfConfig))

	// routes
	routes.RouteSetup(app)

	// port := config.Config("PORT")

	app.Listen(":3000")

	// connection error
	// if conErr != nil {
	// 	panic(conErr)
	// }
}
