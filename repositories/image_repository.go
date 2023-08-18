package repositories

import (
	"clipsearch/models"
	"errors"
)

type ImageRepository interface {
	Count() (int, error)
	CountWithSha256(sha256 string) (int, error)
	// the int is the id of the newly created image
	Create(image *models.ImageModel) (int, error)
	GetImages(offset int, limit int) ([]models.ImageModel, error)
	GetById(id int) (*models.ImageModel, error)
	DeleteById(id int) error
}

var ImageNotFoundError = errors.New("Image with such id was not found")
