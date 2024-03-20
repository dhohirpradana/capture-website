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

	screenshotParamJson, err := json.Marshal(screenshotParam)
	assert.Equal(t, nil, err, err)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(screenshotParamJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	result := screenshot.Capture(c)

	if assert.NoError(t, result) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCaptureZeroTime(t *testing.T) {
	screenshot := InitScreenshot()
	e := echo.New()
	var screenshotParam model.ScreenshotParam

	screenshotParam.Url = "https://www.youtube.com"
	screenshotParam.Filename = "testaja"
	screenshotParam.Wait = 0
	screenshotParam.Quality = 100
	screenshotParam.Width = 1490
	screenshotParam.Height = 1080

	screenshotParamJson, err := json.Marshal(screenshotParam)
	assert.Equal(t, nil, err, err)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(screenshotParamJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	result := screenshot.Capture(c)

	if assert.NoError(t, result) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCaptureNoFilename(t *testing.T) {
	screenshot := InitScreenshot()
	e := echo.New()
	var screenshotParam model.ScreenshotParam

	screenshotParam.Url = "https://www.youtube.com"
	screenshotParam.Wait = 30
	screenshotParam.Quality = 100
	screenshotParam.Width = 1490
	screenshotParam.Height = 1080

	screenshotParamJson, err := json.Marshal(screenshotParam)
	assert.Equal(t, nil, err, err)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(screenshotParamJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	result := screenshot.Capture(c)

	if assert.NoError(t, result) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCaptureInvalidURL(t *testing.T) {
	screenshot := InitScreenshot()
	e := echo.New()
	var screenshotParam model.ScreenshotParam

	screenshotParam.Url = "https://www.youtube.com1"
	screenshotParam.Filename = "testaja"
	screenshotParam.Wait = 5
	screenshotParam.Quality = 100
	screenshotParam.Width = 1490
	screenshotParam.Height = 1080

	screenshotParamJson, err := json.Marshal(screenshotParam)
	assert.Equal(t, nil, err, err)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(screenshotParamJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	result := screenshot.Capture(c)

	if assert.NoError(t, result) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
