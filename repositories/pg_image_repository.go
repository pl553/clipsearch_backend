package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"clipsearch/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgImageRepository struct {
	pool *pgxpool.Pool
}

func NewPgImageRepository(pool *pgxpool.Pool) *PgImageRepository {
	return &PgImageRepository{pool: pool}
}

func (repo *PgImageRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM Images`
	row := repo.pool.QueryRow(context.Background(), query)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("Failed to count images: %w", err)
	}
	return count, nil
}

func (repo *PgImageRepository) CountWithSha256(sha256 string) (int, error) {
	query := `SELECT COUNT(*) FROM Images WHERE Sha256=$1`
	row := repo.pool.QueryRow(context.Background(), query, sha256)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("Failed to count images: %w", err)
	}
	return count, nil
}

func embeddingToString(embedding []float32) string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for i, val := range embedding {
		builder.WriteString(fmt.Sprintf("%f", val))
		if i != len(embedding)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
	return builder.String()
}

func (repo *PgImageRepository) Create(image *models.Image) (int, error) {
	query := `INSERT INTO Images (SourceUrl,ThumbnailUrl,Sha256,Embedding) VALUES ($1,$2,$3,$4) RETURNING ImageID;`
	rows, err := repo.pool.Query(
		context.Background(),
		query, image.SourceUrl,
		image.ThumbnailUrl,
		image.Sha256,
		embeddingToString(image.Embedding))

	if err != nil {
		return 0, fmt.Errorf("Failed to create image: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, errors.New("Failed to get id of newly created image")
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}
	var id int
	err = rows.Scan(&id)
	if err != nil {
		return 0, errors.New("Failed to get id of newly created image")
	}
	return id, nil
}

func (repo *PgImageRepository) GetImages(offset int, limit int) ([]models.Image, error) {
	query := `SELECT ImageID, SourceUrl, ThumbnailUrl, Sha256 FROM Images ORDER BY ImageID LIMIT $1 OFFSET $2;`
	rows, err := repo.pool.Query(context.Background(), query, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to get images: %w", err)
	}
	
	defer rows.Close()

	images := make([]models.Image, 0, 32)

	for rows.Next() {
		var image models.Image
		if err := rows.Scan(&image.ImageID, &image.SourceUrl, &image.ThumbnailUrl, &image.Sha256); err != nil {
			return nil, fmt.Errorf("Failed to get images: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

func (repo *PgImageRepository) GetSimilarImages(embedding []float32, offset int, limit int) ([]models.Image, error) {
	query := `SELECT ImageID, SourceUrl, ThumbnailUrl, Sha256 FROM Images ORDER BY Embedding <#> $1 LIMIT $2 OFFSET $3;`
	rows, err := repo.pool.Query(context.Background(), query, embeddingToString(embedding), limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to get images: %w", err)
	}
	
	defer rows.Close()

	images := make([]models.Image, 0, 32)

	for rows.Next() {
		var image models.Image
		if err := rows.Scan(&image.ImageID, &image.SourceUrl, &image.ThumbnailUrl, &image.Sha256); err != nil {
			return nil, fmt.Errorf("Failed to get images: %w", err)
		}
		images = append(images, image)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return images, nil
}

func (repo *PgImageRepository) GetById(id int) (*models.Image, error) {
	query := "SELECT ImageID,SourceUrl,ThumbnailUrl,Sha256 FROM Images WHERE ImageID=$1"
	rows, err := repo.pool.Query(context.Background(), query, id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get image by id: %w", err)
	}

	defer rows.Close()
	
	if !rows.Next() {
		return nil, ImageNotFoundError
	}
	var image models.Image
	if err := rows.Scan(&image.ImageID, &image.SourceUrl, &image.ThumbnailUrl, &image.Sha256); err != nil {
		return nil, fmt.Errorf("Failed to get image by id: %w", err)
	}
	return &image, nil
}

func (repo *PgImageRepository) DeleteById(id int) error {
	query := "DELETE FROM Images WHERE ImageID=$1"
	commandTag, err := repo.pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ImageNotFoundError
	} else {
		return nil
	}
}
