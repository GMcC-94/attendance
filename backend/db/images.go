package db

import (
	"database/sql"

	"github.com/gmcc94/attendance-go/types"
)

type ImageStore interface {
	SaveImage(img types.Image) error
	GetImageByType(imgType string) (types.Image, error)
}

type PostgresImageStore struct {
	DB *sql.DB
}

func (p *PostgresImageStore) SaveImage(img types.Image) error {
	_, err := p.DB.Exec(`
		INSERT INTO images (file_name, file_url, type, uploaded_at)
		VALUES ($1, $2, $3, $4)
	`, img.FileName, img.FileURL, img.Type, img.UploadedAt)
	return err
}

func (p *PostgresImageStore) GetImageByType(imageType string) (types.Image, error) {
	var img types.Image
	err := p.DB.QueryRow(`
		SELECT id, file_name, file_url, type
		FROM images
		WHERE type = $1
		ORDER BY uploaded_at DESC
		LIMIT 1
	`, imageType).Scan(&img.ID, &img.FileName, &img.FileURL, &img.Type)
	return img, err
}
