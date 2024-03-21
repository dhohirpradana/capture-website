package main

import (
	screenshotHandler "captureWeb/helper"
	"github.com/gofiber/fiber/v2"
)

func main() {
	screenshot := screenshotHandler.InitScreenshot()
	app := fiber.New()

	app.Post("/capture", screenshot.Capture)

	_ = app.Listen(":9090")
}
