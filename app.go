package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	"github.com/raliqala/golang-fibre-boilerplate/src/database"
	"github.com/raliqala/golang-fibre-boilerplate/src/routes"
)

func main() {
	fibreConfig := fiber.Config{
		ServerHeader: config.Config("APP_NAME"),
	}
	app := fiber.New(fibreConfig)

	// default routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "All good here!.. 🚀 ",
		})
	})

	// declaration
	// csrfConfig := csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Strict",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }

	// app config
	app.Use(cors.New())
	// app.Use(csrf.New(csrfConfig))

	// routes
	routes.RouteSetup(app)

	port := config.Config("PORT")

	database.ConnectPg()
	conErr := app.Listen(":" + port)

	// connection error
	if conErr != nil {
		panic(conErr)
	}
}
