package models

import "time"

type ScreenshotParam struct {
	Wait     time.Duration `json:"wait" validate:"nonzero,nonnil"`
	Url      string        `json:"url" validate:"nonzero,nonnil"`
	Width    int64         `json:"width"`
	Height   int64         `json:"height"`
	Quality  int           `json:"quality"`
	Filename string        `json:"filename" validate:"nonzero,nonnil"`
}
