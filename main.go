package main

import (
	screenshotHandler "captureWeb/helper"
	"github.com/gofiber/fiber/v2"
)

func main() {
	screenshot := screenshotHandler.InitScreenshot()
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Capture Website v1.0.1",
	})

	app.Post("/capture", screenshot.Capture)

	_ = app.Listen(":9090")
}
