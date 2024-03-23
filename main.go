package main

import (
	screenshotHandler "captureWeb/helper"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	screenshot := screenshotHandler.InitScreenshot()
	app := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Capture Website v1.0.1",
	})

	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"super": "man",
		},
	}))

	app.Post("/capture", screenshot.Capture)
	app.Get("/metrics", monitor.New())

	_ = app.Listen(":9090")
}
