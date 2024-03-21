package helper

import (
	"captureWeb/entity"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("CAPTURE WEB TEST START")
	m.Run()
	fmt.Println("CAPTURE WEB TEST DONE")
}

func TestCapture(t *testing.T) {
	screenshot := InitScreenshot()
	app := fiber.New()

	app.Post("/capture", screenshot.Capture)

	defer func(app *fiber.App) {
		_ = app.Shutdown()
	}(app)

	tests := []struct {
		route    string
		name     string
		body     entity.ScreenshotParam
		expected int
	}{
		{
			"/not-found",
			"If path not found",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
				Wait:     5,
				Quality:  100,
			},
			http.StatusNotFound,
		},
		{
			"/capture",
			"If body is valid",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
				Wait:     5,
				Quality:  100,
			},
			http.StatusOK,
		},
		{
			"/capture",
			"If body is incomplete",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
			},
			http.StatusUnprocessableEntity,
		},
		{
			"/capture",
			"If url is invalid or timeout",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com1",
				Filename: "filename",
				Quality:  100,
				Wait:     5,
			},
			http.StatusInternalServerError,
		},
		{
			"/capture",
			"If wait is more than 1 minute",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
				Quality:  100,
				Wait:     65,
			},
			http.StatusOK,
		},
		{
			"/capture",
			"If quality more than 100",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
				Quality:  110,
				Wait:     5,
			},
			http.StatusUnprocessableEntity,
		},
		{
			"/capture",
			"If quality less than 1",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
				Quality:  -10,
				Wait:     5,
			},
			http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.body.Width == 0 {
				test.body.Width = 1490
			}

			if test.body.Height == 0 {
				test.body.Height = 1080
			}

			screenshotParamJson, _ := json.Marshal(&test.body)

			req := httptest.NewRequest(fiber.MethodPost, test.route, strings.NewReader(string(screenshotParamJson)))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Add(fiber.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

			resp, err := app.Test(req, -1)
			assert.Nil(t, err)

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			assert.Equalf(t, test.expected, resp.StatusCode, test.name)
		})
	}
}

func BenchmarkCapture(b *testing.B) {
	screenshot := InitScreenshot()
	app := fiber.New()

	app.Post("/capture", screenshot.Capture)

	body := entity.ScreenshotParam{
		Url:      "https://www.youtube.com",
		Filename: "filename",
		Quality:  100,
		Wait:     10,
	}

	for i := 0; i < b.N; i++ {
		if body.Width == 0 {
			body.Width = 1490
		}

		if body.Height == 0 {
			body.Height = 1080
		}

		screenshotParamJson, _ := json.Marshal(body)

		req := httptest.NewRequest(fiber.MethodPost, "/capture", strings.NewReader(string(screenshotParamJson)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		req.Header.Add(fiber.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

		_, _ = app.Test(req, -1)
	}
}

func BenchmarkTest(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		j += i
	}
}
