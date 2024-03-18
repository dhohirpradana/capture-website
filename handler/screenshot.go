package handler

import (
	"captureWeb/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/signintech/gopdf"
	"gopkg.in/validator.v2"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ScreenshotHandler struct {
}

func InitScreenshot() ScreenshotHandler {
	return ScreenshotHandler{}
}

func fullScreenshot(waitSec time.Duration, url string, quality int, width int64, height int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(width, height),
		chromedp.Navigate(url),
		chromedp.Sleep(waitSec * time.Second),
		chromedp.FullScreenshot(res, quality),
	}
}

func deleteFile(uuid uuid.UUID) {
	pngPath := "result/capture-" + uuid.String() + ".png"
	pdfPath := "result/capture-" + uuid.String() + ".pdf"

	// delete png
	pngErr := os.Remove(pngPath)
	if pngErr != nil {
		fmt.Println("Error: ", pngErr)
	} else {
		fmt.Println("Successfully deleted file: ", pngPath)
	}

	// delete pdf
	pdfErr := os.Remove(pdfPath)
	if pdfErr != nil {
		fmt.Println("Error: ", pdfErr)
	} else {
		fmt.Println("Successfully deleted file: ", pdfPath)
	}
}

func (h ScreenshotHandler) Capture(c echo.Context) (err error) {
	body, err := io.ReadAll(c.Request().Body)

	id := uuid.New()
	filePath := "result/capture-" + id.String()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var screenshotParam models.ScreenshotParam
	err = json.Unmarshal(body, &screenshotParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := validator.Validate(screenshotParam); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var buf []byte

	url := screenshotParam.Url
	wait := &screenshotParam.Wait
	width := &screenshotParam.Width
	height := &screenshotParam.Height
	quality := &screenshotParam.Quality
	filename := &screenshotParam.Filename

	if *quality == 0 {
		*quality = 100
	}

	if *wait == 0 {
		*wait = 1
	}

	if *width == 0 {
		*width = 1490
	}

	if *height == 0 {
		*height = 1080
	}

	fmt.Println("screenShotParam:", screenshotParam)

	if err := chromedp.Run(ctx, fullScreenshot(*wait, url, *quality, *width, *height, &buf)); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filePath+".png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA3})
	pdf.AddPage()

	err = pdf.Image(filePath+".png", 0, 0, nil)

	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.WritePdf(filePath + ".pdf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	filenameNew := *filename + ".pdf"

	// delete files
	defer deleteFile(id)

	return c.Attachment(filePath+".pdf", filenameNew)
}
