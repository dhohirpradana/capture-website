package handler

import (
	"captureWeb/model"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	e := echo.New()
	var screenshotParam model.ScreenshotParam

	screenshotParam.Url = "https://www.youtube.com"
	screenshotParam.Filename = "testaja"
	screenshotParam.Wait = 5
	screenshotParam.Quality = 100
	screenshotParam.Width = 1490
	screenshotParam.Height = 1080

	tests := []struct {
		name     string
		body     model.ScreenshotParam
		expected int
	}{
		{
			"Valid Body",
			model.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "customfilename",
				Wait:     5,
			},
			http.StatusOK,
		},
		{
			"Invalid Body",
			model.ScreenshotParam{
				Url:      "https://www.youtube.com",
				Filename: "customfilename",
			},
			http.StatusUnprocessableEntity,
		},
		{
			"Invalid URL Timeout",
			model.ScreenshotParam{
				Url:      "https://www.youtube.com1",
				Filename: "customfilename",
				Wait:     5,
			},
			http.StatusGatewayTimeout,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.body.Quality == 0 {
				test.body.Quality = 100
			}

			if test.body.Wait == 0 {
				test.body.Wait = 1
			}

			if test.body.Width == 0 {
				test.body.Width = 1490
			}

			if test.body.Height == 0 {
				test.body.Height = 1080
			}

			screenshotParamJson, err := json.Marshal(test.body)
			assert.Equal(t, nil, err, err)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(screenshotParamJson)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			result := screenshot.Capture(c)

			if assert.NoError(t, result) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}
