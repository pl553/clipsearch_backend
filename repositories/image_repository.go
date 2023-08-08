package repositories

import "clipsearch/models"

type ImageRepository interface {
	Count() (int, error)
	Create(image *models.ImageModel) error
	GetImages(offset int, limit int) ([]models.ImageModel, error)
	GetById(id int) (*models.ImageModel, error)
}
