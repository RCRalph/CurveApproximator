package controllers

import "github.com/gofiber/fiber/v2"

func Index(c *fiber.Ctx) error {
	return c.SendFile("client/dist/index.html")
}
