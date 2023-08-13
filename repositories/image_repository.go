package repositories

import "clipsearch/models"

type ImageRepository interface {
	Count() (int, error)
	// the int is the id of the newly created image
	Create(image *models.ImageModel) (int, error)
	GetImages(offset int, limit int) ([]models.ImageModel, error)
	GetById(id int) (*models.ImageModel, error)
}
