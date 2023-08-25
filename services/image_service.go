package services

import (
	"bytes"
	"clipsearch/config"
	"clipsearch/models"
	"clipsearch/repositories"
	"clipsearch/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type ImageService struct {
	ImageRepo repositories.ImageRepository
	clip      ClipService
}

func NewImageService(imageRepo repositories.ImageRepository, clipService ClipService) *ImageService {
	return &ImageService{
		ImageRepo: imageRepo,
		clip:      clipService,
	}
}

func (s *ImageService) GetCountAndImages(offset int, limit int) (int, []models.ImageModel, error) {
	count, err := s.ImageRepo.Count()
	if err != nil {
		return 0, nil, err
	}

	images, err := s.ImageRepo.GetImages(offset, limit)
	if err != nil {
		return 0, nil, err
	}

	return count, images, nil
}

var ImageExistsError = fmt.Errorf("This image already exists (hash match)")

func (s *ImageService) AddImageByURL(url string, thumbnailUrl string) error {
	var buf bytes.Buffer

	err := utils.DownloadFile(&buf, url, config.MAX_IMAGE_FILE_SIZE)
	if err != nil {
		return err
	}

	hashSHA256 := sha256.New()
	if _, err := io.Copy(hashSHA256, bytes.NewReader(buf.Bytes())); err != nil {
		return err
	}
	hashBytes := hashSHA256.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	count, err := s.ImageRepo.CountWithSha256(hashString)
	if err != nil {
		return err
	}
	if count > 0 {
		return ImageExistsError
	}

	embedding, err := s.clip.EncodeImage(buf.Bytes())
	if err != nil {
		return err
	}

	image := models.ImageModel{
		SourceUrl:    url,
		ThumbnailUrl: thumbnailUrl,
		Sha256:       hashString,
		Embedding:    embedding,
	}

	_, err = s.ImageRepo.Create(&image)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImageService) GetImagesSimilarToText(textPrompt string, offset int, limit int) ([]models.ImageModel, error) {
	textEmbedding, err := s.clip.EncodeText(textPrompt)

	if err != nil {
		return nil, err
	}

	return s.ImageRepo.GetSimilarImages(textEmbedding, offset, limit)
}
