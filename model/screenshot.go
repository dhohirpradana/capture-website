package model

import "time"

type ScreenshotParam struct {
	Wait     time.Duration `json:"wait" query:"wait" validate:"nonzero,nonnil,min=1"`
	Url      string        `json:"url" query:"url" validate:"nonzero,nonnil"`
	Width    int64         `json:"width" query:"width"`
	Height   int64         `json:"height"  query:"height"`
	Quality  int           `json:"quality" query:"quality" validate:"nonzero,nonnil,min=1,max=100"`
	Filename string        `json:"filename" query:"filename" validate:"nonzero,nonnil"`
}
