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

func (repo *MockImageRepository) GetSimilarImages(embedding []float32, offset int, limit int) ([]models.ImageModel, error) {
	return nil, nil
}

func (repo *MockImageRepository) CountWithSha256(sha256 string) (int, error) {
	counter := 0
	for _, image := range repo.images {
		if image.Sha256 == sha256 {
			counter++
		}
	}
	return counter, nil
}

func (repo *MockImageRepository) Create(image *models.ImageModel) (int, error) {
	newImage := *image
	newImage.ImageID = repo.ct + 1
	repo.ct++
	repo.images = append(repo.images, newImage)
	return newImage.ImageID, nil
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
	return nil, ImageNotFoundError
}

func (repo *MockImageRepository) DeleteById(id int) error {
	for i, image := range repo.images {
		if image.ImageID == id {
			repo.images[i] = repo.images[len(repo.images)-1]
			repo.images = repo.images[:len(repo.images)-1]
			return nil
		}
	}
	return ImageNotFoundError
}
