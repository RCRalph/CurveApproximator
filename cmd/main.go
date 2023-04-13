package main

import (
	"log"

	"github.com/RCRalph/CurveApproximator/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/", "client/dist/public")

	app.Get("/", controllers.Index)

	log.Fatal(app.Listen(":8080"))
}
