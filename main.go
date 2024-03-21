package main

import (
	screenshotHandler "captureWeb/helper"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	screenshot := screenshotHandler.InitScreenshot()
	e := echo.New()

	e.POST("/capture", screenshot.Capture)

	// Start the server
	err := e.Start(":9090")
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
