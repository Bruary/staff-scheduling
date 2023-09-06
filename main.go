package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/user", func(c *fiber.Ctx) error {
		return c.SendString("I got the user!")
	})

	app.Listen(":3000")
}
