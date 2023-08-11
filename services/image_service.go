package services

import (
	"bytes"
	"clipsearch/config"
	"clipsearch/models"
	"clipsearch/repositories"
	"clipsearch/utils"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type ImageService struct {
	imageRepo repositories.ImageRepository
}

func NewImageService(imageRepo repositories.ImageRepository) *ImageService {
	return &ImageService{imageRepo: imageRepo}
}

func (s *ImageService) GetCountAndImages(offset int, limit int) (int, []models.ImageModel, error) {
	count, err := s.imageRepo.Count()
	if err != nil {
		return 0, nil, err
	}

	images, err := s.imageRepo.GetImages(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	return count, images, nil
}

func (s *ImageService) AddImageByURL(url string) error {
	var buf bytes.Buffer

	err := utils.DownloadFile(&buf, url, config.MAX_IMAGE_FILE_SIZE)
	if err != nil {
		return err
	}

	hashSHA256 := sha256.New()
	if _, err := io.Copy(hashSHA256, &buf); err != nil {
		return err
	}
	hashBytes := hashSHA256.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	image := models.ImageModel{
		SourceUrl: url,
		Sha256:    hashString,
	}

	if err := s.imageRepo.Create(&image); err != nil {
		return err
	}

	return nil
}
