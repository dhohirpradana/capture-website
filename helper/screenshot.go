package helper

import (
	"captureWeb/entity"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"github.com/signintech/gopdf"
	"gopkg.in/validator.v2"
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

func (h ScreenshotHandler) Capture(c *fiber.Ctx) (err error) {
	tempDir, err := createTempDir()
	fmt.Println(tempDir)

	defer func(tempDir string) {
		err := removeTempDir(tempDir)
		if err != nil {
			fmt.Println(err.Error())
		}
	}(tempDir)

	var screenshotParam entity.ScreenshotParam

	if err := c.BodyParser(&screenshotParam); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := validator.Validate(screenshotParam); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
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

	if err := chromedp.Run(ctx, fullScreenshot(*wait, url, *quality, *width, *height, &buf)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	tempPngPath := filepath.Join(tempDir, *filename+".png")
	err = writeFile(tempPngPath, buf)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA3})
	pdf.AddPage()

	err = pdf.Image(tempPngPath, 0, 0, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	tempPdfPath := filepath.Join(tempDir, *filename+".pdf")
	err = pdf.WritePdf(tempPdfPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendFile(tempPdfPath)
}
