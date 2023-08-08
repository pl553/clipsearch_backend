package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"clipsearch/binding"
	"clipsearch/config"
	"clipsearch/models"
	"clipsearch/repositories"
	"clipsearch/utils"

	"github.com/clevergo/jsend"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var imageRepository *repositories.ImageRepository
var internalErrorJson = jsend.NewError("Internal error", 500, nil)

type PostImagesForm struct {
	Url string `schema:"url,required" validate:"url"`
}

func postImages(c *gin.Context) {
	c.Request.ParseMultipartForm(2048)
	var form PostImagesForm
	if err := binding.ShouldBind(&form, c.Request.Form); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}

	imagePath, err := utils.DownloadFile(form.Url, config.MAX_IMAGE_FILE_SIZE)
	if err != nil {
		if errors.Is(err, utils.FileSizeExceededError{}) {
			c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{
				"url": fmt.Sprintf("Image at url is too large (>%d MB)", config.MAX_IMAGE_FILE_SIZE_MB),
			}))
			return
		} else {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, internalErrorJson)
			return
		}
	}

	defer os.Remove(imagePath)
	imageFile, err := os.Open(imagePath)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}
	hashSHA256 := sha256.New()
	if _, err := io.Copy(hashSHA256, imageFile); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}
	hashBytes := hashSHA256.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	image := models.ImageModel{
		SourceUrl: form.Url,
		Sha256:    hashString,
	}

	if err := imageRepository.Create(&image); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}
	c.JSON(http.StatusOK, jsend.New(nil))
}

type GetImagesQuery struct {
	Offset int `schema:"offset" validate:"min=0"`
	Limit  int `schema:"limit" validate:"min=0"`
}

func getImages(c *gin.Context) {
	var query GetImagesQuery

	if err := binding.ShouldBind(&query, c.Request.URL.Query()); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}

	count, err := imageRepository.Count()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	images, err := imageRepository.GetImages(query.Offset, query.Limit)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	c.JSON(http.StatusOK, jsend.New(gin.H{
		"imageCount": count,
		"images":     images,
	}))
}

func ginParamsToMap(params gin.Params) map[string][]string {
	result := make(map[string][]string)
	for _, kv := range params {
		result[kv.Key] = []string{kv.Value}
	}
	return result
}

type GetImageByIdQuery struct {
	Id int `schema:"id" validate:"min=0"`
}

func getImageById(c *gin.Context) {
	var query GetImageByIdQuery
	if err := binding.ShouldBind(&query, ginParamsToMap(c.Params)); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}
	image, err := imageRepository.GetById(query.Id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}
	if image == nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{"id": "No image with such id exists"}))
		return
	}
	c.JSON(http.StatusOK, jsend.New(image))
}

func main() {
	port := os.Getenv(config.PORT_ENVAR)
	if port == "" {
		port = config.DEFAULT_PORT
	}
	dbConnString := os.Getenv(config.DATABASE_CONNECTION_URL_ENVAR)
	if dbConnString == "" {
		log.Fatal(fmt.Sprintf("Please define the %v envar", config.DATABASE_CONNECTION_URL_ENVAR))
	}
	pgPool, err := pgxpool.New(context.Background(), dbConnString)
	if err != nil {
		log.Print("Failed to connect to db!")
		log.Fatal(err)
	}
	defer pgPool.Close()

	imageRepository = repositories.NewImageRepository(pgPool)

	count, err := imageRepository.Count()
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		seedImages := []models.ImageModel{
			{SourceUrl: "/static/images/1.gif"},
			{SourceUrl: "/static/images/2.jpg"},
			{SourceUrl: "/static/images/3.jpg"},
		}

		for _, image := range seedImages {
			if err := imageRepository.Create(&image); err != nil {
				log.Fatal(err)
			}
		}
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/api/images", getImages)
	router.POST("/api/images", postImages)
	router.GET("/api/images/:id", getImageById)
	router.Run("localhost:" + port)
}
