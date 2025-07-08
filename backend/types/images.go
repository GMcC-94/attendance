package types

import "time"

type Image struct {
	ID         int       `json:"id"`
	FileName   string    `json:"fileName"`
	FilePath   string    `json:"filePath"`
	Type       string    `json:"type"`
	UploadedAt time.Time `json:"uploadedAt"`
}
