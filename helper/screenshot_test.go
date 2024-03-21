package helper

import (
	"captureWeb/entity"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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

	defer app.Shutdown()

	tests := []struct {
		name     string
		body     entity.ScreenshotParam
		expected int
	}{
		{
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
			"If body is incomplete",
			entity.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "filename",
			},
			http.StatusUnprocessableEntity,
		},
		{
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

			screenshotParamJson, err := json.Marshal(test.body)
			assert.Equal(t, nil, err, err)

			req := httptest.NewRequest(fiber.MethodPost, "/capture", strings.NewReader(string(screenshotParamJson)))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Add(fiber.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

			resp, _ := app.Test(req, 1)

			assert.Equalf(t, test.expected, resp.StatusCode, test.name)
		})
	}
}

func BenchmarkCapture(b *testing.B) {
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

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(screenshotParamJson)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		_ = httptest.NewRecorder()
	}
}

func BenchmarkTest(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		j += i
	}
}
