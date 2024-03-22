package helper

import (
	"captureWeb/entity"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/signintech/gopdf"
	"gopkg.in/validator.v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type ScreenshotHandler struct {
}

func InitScreenshot() ScreenshotHandler {
	return ScreenshotHandler{}
}

func fullScreenshot(screenshotParam *entity.ScreenshotParam, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(screenshotParam.Width, screenshotParam.Height),
		chromedp.Navigate(screenshotParam.Url),
		chromedp.Sleep(screenshotParam.Wait * time.Second),
		chromedp.FullScreenshot(res, screenshotParam.Quality),
	}
}

func writeFile(tempPath string, buf []byte) error {
	if err := os.WriteFile(tempPath, buf, 0o666); err != nil {
		return err
	}
	return nil
}

func createTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "dir")
	if err != nil {
		return tempDir, err
	}
	return tempDir, nil
}

func removeTempDir(tempDir string) error {
	err := os.RemoveAll(tempDir)
	if err != nil {
		return err
	}
	return nil
}

func (h ScreenshotHandler) Capture(c echo.Context) (err error) {
	tempDir, err := createTempDir()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer func(path string) {
		err := removeTempDir(path)
		if err != nil {
			fmt.Println(err.Error())
		}
	}(tempDir)

	var screenshotParam entity.ScreenshotParam

	err = c.Bind(&screenshotParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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

	urlBody := &screenshotParam.Url
	width := &screenshotParam.Width
	height := &screenshotParam.Height
	filename := &screenshotParam.Filename

	_, err = url.ParseRequestURI(*urlBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if *width == 0 {
		*width = 1490
	}

	if *height == 0 {
		*height = 1080
	}

	if err := chromedp.Run(ctx, fullScreenshot(&screenshotParam, &buf)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	tempPngPath := filepath.Join(tempDir, *filename+".png")
	if err := writeFile(tempPngPath, buf); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA3})
	pdf.AddPage()

	err = pdf.Image(tempPngPath, 0, 0, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	tempPdfPath := filepath.Join(tempDir, *filename+".pdf")
	err = pdf.WritePdf(tempPdfPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Attachment(tempPdfPath, *filename+".pdf")
}
