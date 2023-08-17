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
	ImageRepo                   repositories.ImageRepository
	imageFeatureExtractionQueue chan models.ImageModel
}

func NewImageService(ImageRepo repositories.ImageRepository) *ImageService {
	return &ImageService{
		ImageRepo:                   ImageRepo,
		imageFeatureExtractionQueue: make(chan models.ImageModel, 1024),
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

	image := models.ImageModel{
		SourceUrl:    url,
		ThumbnailUrl: thumbnailUrl,
		Sha256:       hashString,
	}

	count, err := s.ImageRepo.CountWithSha256(image.Sha256)
	if err != nil {
		return err
	}
	if count > 0 {
		return ImageExistsError
	}

	id, err := s.ImageRepo.Create(&image)
	if err != nil {
		return err
	}
	image.ImageID = id

	conn, err := ConnectToImageFeatureDaemon("tcp://localhost:" + config.FEATURE_EXTRACT_DAEMON_PORT)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.SendImageForFeatureExtraction(id, buf.Bytes()); err != nil {
		return err
	}

	return nil
}
