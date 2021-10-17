package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/raliqala/safepass_api/src/config"
	Database "github.com/raliqala/safepass_api/src/db"
	"github.com/raliqala/safepass_api/src/routes"
)

func main() {
	fibreConfig := fiber.Config{
		ServerHeader: config.Config("APP_NAME"),
	}
	app := fiber.New(fibreConfig)

	// declaration
	csrfConfig := csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
	}

	// app config
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(csrf.New(csrfConfig))

	// routes
	routes.RouteSetup(app)

	port := config.Config("PORT")

	Database.ConnectPg()
	conErr := app.Listen(":" + port)

	// connection error
	if conErr != nil {
		panic(conErr)
	}
}
