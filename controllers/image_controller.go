package controllers

import (
	"clipsearch/binding"
	"clipsearch/config"
	"clipsearch/services"
	"clipsearch/utils"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/clevergo/jsend"
	"github.com/gin-gonic/gin"
)

var internalErrorJson = jsend.NewError("Internal error", 500, nil)

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController(imageService *services.ImageService) *ImageController {
	return &ImageController{imageService: imageService}
}

type PostImagesForm struct {
	Url string `schema:"url,required" validate:"url"`
}

func (controller *ImageController) PostImages(c *gin.Context) {
	c.Request.ParseMultipartForm(2048)
	var form PostImagesForm
	if err := binding.ShouldBind(&form, c.Request.Form); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}

	if err := controller.imageService.AddImageByURL(form.Url); err != nil {
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

	c.JSON(http.StatusOK, jsend.New(nil))
}

type GetImagesQuery struct {
	Offset int `schema:"offset" validate:"min=0"`
	Limit  int `schema:"limit" validate:"min=0"`
}

func (controller *ImageController) GetImages(c *gin.Context) {
	var query GetImagesQuery

	if err := binding.ShouldBind(&query, c.Request.URL.Query()); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}

	count, images, err := controller.imageService.GetCountAndImages(query.Offset, query.Limit)
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

func (controller *ImageController) GetImageById(c *gin.Context) {
	var query GetImageByIdQuery
	if err := binding.ShouldBind(&query, ginParamsToMap(c.Params)); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}
	image, err := controller.imageService.ImageRepo.GetById(query.Id)
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
