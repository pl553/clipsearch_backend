package repositories

import (
	"clipsearch/models"
)

type MockImageRepository struct {
	images []models.ImageModel
	ct     int
}

func NewMockImageRepository() *MockImageRepository {
	return &MockImageRepository{images: make([]models.ImageModel, 0, 16), ct: 0}
}

func (repo *MockImageRepository) Count() (int, error) {
	return len(repo.images), nil
}

func (repo *MockImageRepository) Create(image *models.ImageModel) error {
	newImage := *image
	newImage.ImageID = repo.ct + 1
	repo.ct++
	repo.images = append(repo.images, newImage)
	return nil
}

func (repo *MockImageRepository) GetImages(offset int, limit int) ([]models.ImageModel, error) {
	return repo.images[offset : offset+limit], nil
}

func (repo *MockImageRepository) GetById(id int) (*models.ImageModel, error) {
	for _, image := range repo.images {
		if image.ImageID == id {
			return &image, nil
		}
	}
	return nil, nil
}
