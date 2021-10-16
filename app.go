package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/config"
	"github.com/raliqala/safepass_api/src/db"

	"log"
)

func main() {
	app := fiber.New()

	db.ConnectPg()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/post", func(c *fiber.Ctx) error {
		return c.SendString("New post posted 📮 ")
	})

	log.Fatal(app.Listen(config.Config(":PORT")))
}
