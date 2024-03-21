package helper

import (
	"captureWeb/entity"
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
	"path/filepath"
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

func (h ScreenshotHandler) Capture(c echo.Context) (err error) {
	id := uuid.New()
	filePath := "capture-" + id.String()

	tempDir, err := os.MkdirTemp("", "dir")
	if err != nil {
		log.Fatal(err)
	}

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println(err.Error())
		}
	}(tempDir)

	fmt.Println(tempDir)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var screenshotParam entity.ScreenshotParam

	body, err := io.ReadAll(c.Request().Body)
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

	if *width == 0 {
		*width = 1490
	}

	if *height == 0 {
		*height = 1080
	}

	//fmt.Println("screenShotParam:", screenshotParam)

	if err := chromedp.Run(ctx, fullScreenshot(*wait, url, *quality, *width, *height, &buf)); err != nil {
		//log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	tempPngPath := filepath.Join(tempDir, filePath+".png")
	if err := os.WriteFile(tempPngPath, buf, 0o644); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA3})
	pdf.AddPage()

	err = pdf.Image(tempPngPath, 0, 0, nil)

	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	tempPdfPath := filepath.Join(tempDir, filePath+".pdf")
	err = pdf.WritePdf(tempPdfPath)
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	filenameNew := *filename + ".pdf"

	return c.Attachment(tempPdfPath, filenameNew)
}
