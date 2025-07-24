package types

import "time"

type Image struct {
	ID         int       `json:"id"`
	FileName   string    `json:"fileName"`
	FileURL    string    `json:"fileUrl"`
	Type       string    `json:"type"`
	UploadedAt time.Time `json:"uploadedAt"`
}
