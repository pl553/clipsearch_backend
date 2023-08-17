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
	Url          string `schema:"url,required" validate:"url"`
	ThumbnailUrl string `schema:"thumbnailUrl" validate:"omitempty,url"`
}

func (controller *ImageController) PostImages(c *gin.Context) {
	c.Request.ParseMultipartForm(2048)
	var form PostImagesForm
	if err := binding.ShouldBind(&form, c.Request.Form); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}

	if form.ThumbnailUrl == "" {
		form.ThumbnailUrl = form.Url
	}

	if err := controller.imageService.AddImageByURL(form.Url, form.ThumbnailUrl); err != nil {
		if errors.Is(err, utils.FileSizeExceededError{}) {
			c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{
				"url": fmt.Sprintf("Image at url is too large (>%d MB)", config.MAX_IMAGE_FILE_SIZE_MB),
			}))
			return
		} else if errors.Is(err, services.ImageExistsError) {
			c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{
				"url": fmt.Sprintf("Image already exists"),
			}))
			log.Print("attempted to add already existing image")
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
		"totalCount": count,
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
		c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{"id": "=No image with such id exists"}))
		return
	}
	c.JSON(http.StatusOK, jsend.New(image))
}

type SearchQuery struct {
	Query  string `schema:"q"`
	Offset int    `schema:"offset" validate:"min=0"`
	Limit  int    `schema:"limit" validate:"min=0"`
}

func (controller *ImageController) GetSearchImages(c *gin.Context) {
	var query SearchQuery
	if err := binding.ShouldBind(&query, c.Request.URL.Query()); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err.(binding.BindingError).FieldErrors))
		return
	}

	conn, err := services.ConnectToImageSearchDaemon("tcp://localhost:" + config.SEARCH_DAEMON_PORT)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	var daemonResults []services.ImageSearchResult
	var count int

	if len(query.Query) == 0 {
		daemonResults = nil
		count = 0
	} else {
		count, daemonResults, err = conn.SearchRequest(query.Query)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, internalErrorJson)
			return
		}
	}

	type ImageSearchResult struct {
		ImageID      int     `json:"id"`
		ThumbnailUrl string  `json:"thumbnailUrl"`
		Score        float32 `json:"score"`
	}

	results := make([]ImageSearchResult, 0, query.Limit)
	for i := query.Offset; i < query.Offset+query.Limit && i < len(daemonResults); i++ {
		image, err := controller.imageService.ImageRepo.GetById(daemonResults[i].ImageID)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, internalErrorJson)
			return
		}
		results = append(results, ImageSearchResult{
			ImageID:      image.ImageID,
			ThumbnailUrl: image.ThumbnailUrl,
			Score:        daemonResults[i].Score,
		})
	}
	c.JSON(http.StatusOK, jsend.New(gin.H{
		"totalCount": count,
		"images":     results,
	}))
}
