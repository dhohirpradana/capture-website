package helper

import (
	"captureWeb/entity"
	"context"
	"encoding/json"
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
	pngPath := "capture-" + uuid.String() + ".png"
	pdfPath := "capture-" + uuid.String() + ".pdf"

	// delete png
	_ = os.Remove(pngPath)

	// delete pdf
	_ = os.Remove(pdfPath)
}

func (h ScreenshotHandler) Capture(c echo.Context) (err error) {
	id := uuid.New()
	filePath := "capture-" + id.String()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var screenshotParam entity.ScreenshotParam

	if c.Request().Method == "POST" {
		body, err := io.ReadAll(c.Request().Body)
		err = json.Unmarshal(body, &screenshotParam)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := c.Bind(&screenshotParam); err != nil {
			return err
		}
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
