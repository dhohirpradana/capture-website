package main

import (
	screenshotHandler "captureWeb/handler"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	screenshot := screenshotHandler.InitScreenshot()
	echoServer := echo.New()

	echoServer.POST("/capture", screenshot.Capture)

	// Start the server
	err := echoServer.Start(":9090")
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
