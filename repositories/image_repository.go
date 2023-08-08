package repositories

import (
	"context"
	"fmt"

	"clipsearch/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{pool: pool}
}

func (repo *ImageRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM Images`
	row := repo.pool.QueryRow(context.Background(), query)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("Failed to count images: %w", err)
	}
	return count, nil
}

func (repo *ImageRepository) Create(image *models.ImageModel) error {
	query := `INSERT INTO Images (SourceUrl,Sha256) VALUES ($1,$2);`
	rows, err := repo.pool.Query(context.Background(), query, image.SourceUrl, image.Sha256)
	rows.Close()
	if err != nil {
		return fmt.Errorf("Failed to create image: %w", err)
	}
	return nil
}

func (repo *ImageRepository) GetImages(offset int, limit int) ([]models.ImageModel, error) {
	query := `SELECT ImageID, SourceUrl FROM Images ORDER BY ImageID LIMIT $1 OFFSET $2;`
	rows, err := repo.pool.Query(context.Background(), query, limit, offset)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to get images: %w", err)
	}

	images := make([]models.ImageModel, 0, 32)

	for rows.Next() {
		var image models.ImageModel
		if err := rows.Scan(&image.ImageID, &image.SourceUrl); err != nil {
			return nil, fmt.Errorf("Failed to get images: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

func (repo *ImageRepository) GetById(id int) (*models.ImageModel, error) {
	query := "SELECT ImageID,SourceUrl,Sha256 FROM Images WHERE ImageID=$1"
	rows, err := repo.pool.Query(context.Background(), query, id)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to get image by id: %w", err)
	}
	if !rows.Next() {
		return nil, nil
	}
	var image models.ImageModel
	if err := rows.Scan(&image.ImageID, &image.SourceUrl, &image.Sha256); err != nil {
		return nil, fmt.Errorf("Failed to get image by id: %w", err)
	}
	return &image, nil
}
