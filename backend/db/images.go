package db

import (
	"database/sql"

	"github.com/gmcc94/attendance-go/types"
)

type ImageStore interface {
	SaveImage(img types.Image) error
	GetLatestLogo() (types.Image, error)
}

type PostgresImageStore struct {
	DB *sql.DB
}

func (p *PostgresImageStore) SaveImage(img types.Image) error {
	_, err := p.DB.Exec(`
		INSERT INTO images (file_name, file_path, type)
		VALUES ($1, $2, $3)
	`, img.FileName, img.FilePath, img.Type)
	return err
}

func (p *PostgresImageStore) GetLatestLogo() (types.Image, error) {
	var img types.Image
	err := p.DB.QueryRow(`
		SELECT id, file_name, file_path, type, uploaded_at
		FROM images
		WHERE type = 'logo'
		ORDER BY uploaded_at DESC
		LIMIT 1
	`).Scan(&img.ID, &img.FileName, &img.FilePath, &img.Type, &img.UploadedAt)
	return img, err
}
