package helper

import (
	"captureWeb/entity"
	"context"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/signintech/gopdf"
	"gopkg.in/validator.v2"
	"log"
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
	pngPath := "capture-" + uuid.String() + ".png"
	pdfPath := "capture-" + uuid.String() + ".pdf"

	// delete png
	_ = os.Remove(pngPath)

	// delete pdf
	_ = os.Remove(pdfPath)
}

func (h ScreenshotHandler) Capture(c *fiber.Ctx) (err error) {
	id := uuid.New()
	filePath := "capture-" + id.String()

	var screenshotParam entity.ScreenshotParam

	if err := c.BodyParser(screenshotParam); err != nil {
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
	//filename := &screenshotParam.Filename

	if *width == 0 {
		*width = 1490
	}

	if *height == 0 {
		*height = 1080
	}

	//fmt.Println("screenShotParam:", screenshotParam)

	if err := chromedp.Run(ctx, fullScreenshot(*wait, url, *quality, *width, *height, &buf)); err != nil {
		//log.Fatal(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

	//filenameNew := *filename + ".pdf"

	// delete files
	defer deleteFile(id)

	return c.SendFile(filePath + ".pdf")
}
