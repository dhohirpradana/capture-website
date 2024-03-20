package model

import "time"

type ScreenshotParam struct {
	Wait     time.Duration `json:"wait" query:"wait" validate:"nonzero,nonnil"`
	Url      string        `json:"url" query:"url" validate:"nonzero,nonnil"`
	Width    int64         `json:"width" query:"width"`
	Height   int64         `json:"height"  query:"height"`
	Quality  int           `json:"quality" query:"quality"`
	Filename string        `json:"filename" query:"filename" validate:"nonzero,nonnil"`
}
