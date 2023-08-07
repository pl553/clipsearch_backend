package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"clipsearch/binding"
	"clipsearch/config"
	"clipsearch/models"
	"clipsearch/repositories"

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

	url, err := url.ParseRequestURI(form.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{"url": "Invalid url"}))
		return
	}
	if len(url.Query()) != 0 {
		c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{"url": "Url must not have any query parameters in it"}))
		return
	}

	// TODO add more validation (check that its actually an image, check that the filesize isnt too large, ...)

	image := models.ImageModel{SourceUrl: form.Url}
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
		"image_count": count,
		"images":      images,
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
		c.JSON(http.StatusInternalServerError, jsend.NewError("Internal error", 500, nil))
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
