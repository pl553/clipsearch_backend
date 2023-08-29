package controllers

import (
	"clipsearch/binding"
	"clipsearch/config"
	"clipsearch/dtos"
	"clipsearch/repositories"
	"clipsearch/services"
	"clipsearch/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var internalErrorJson = dtos.NewJsendErrorResponse("Internal error")

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController(imageService *services.ImageService) *ImageController {
	return &ImageController{imageService: imageService}
}

// @Summary Get images
// @Description Returns an array of images from the repository, ordered by ID, skipping the first `offset` images and returning at most `limit`.
// @Tags images
// @Produce json
// @Param offset query int false "How many images to skip"
// @Param limit query int false "How many images to return at most"
// @Success 200 {object} dtos.JsendImagesResponse "Success"
// @Failure 400 {object} dtos.JsendFailResponse "Failure (bad params)"
// @Failure 500 {object} dtos.JsendErrorResponse "Failure (internal error)"
// @Router /api/images [get]
func (controller *ImageController) GetImages(c *gin.Context) {
	var query GetImagesQuery

	if err := binding.ShouldBind(&query, c.Request.URL.Query()); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(err.(binding.BindingError).FieldErrors))
		return
	}

	count, images, err := controller.imageService.GetCountAndImages(query.Offset, query.Limit)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	c.JSON(http.StatusOK, dtos.NewJsendImagesResponse(count, images))
}

type SearchQuery struct {
	Query  string `schema:"q"`
	Offset int    `schema:"offset" validate:"min=0"`
	Limit  int    `schema:"limit" validate:"min=0"`
}

// @Summary Search the image repository (text query)
// @Description Returns an array of images from the repository, ordered by relevance, skipping the first `offset` images and returning at most `limit`.
// @Tags search
// @Produce json
// @Param q query string true "The text query"
// @Param offset query int false "How many images to skip"
// @Param limit query int false "How many images to return at most"
// @Success 200 {object} dtos.JsendImagesResponse "Success"
// @Failure 400 {object} dtos.JsendFailResponse "Failure (bad params)"
// @Failure 500 {object} dtos.JsendErrorResponse "Failure (internal error)"
// @Router /api/images/search [get]
func (controller *ImageController) GetSearchImages(c *gin.Context) {
	var query SearchQuery
	if err := binding.ShouldBind(&query, c.Request.URL.Query()); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(err.(binding.BindingError).FieldErrors))
		return
	}

	count, err := controller.imageService.ImageRepo.Count()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	results, err := controller.imageService.GetImagesSimilarToText(query.Query, query.Offset, query.Limit)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	c.JSON(http.StatusOK, dtos.NewJsendImagesResponse(count, results))
}

type PostImagesForm struct {
	Url          string `schema:"url,required" validate:"url"`
	ThumbnailUrl string `schema:"thumbnailUrl" validate:"omitempty,url"`
}

// @Summary Create image
// @Description Adds an image to the repository.
// @Description Image is not added if it already exists in the repository (hash match), or if the file size is larger than allowed (see config)
// @Tags images
// @Produce json
// @Param sourceUrl formData string true "URL of the image to be added."
// @Param thumbnailUrl formData string false "URL to store as thumbnail for the image. Default is source URL."
// @Success 200 {object} dtos.JsendEmptySuccessResponse "Success"
// @Failure 400 {object} dtos.JsendFailResponse "Failure (bad params)"
// @Failure 500 {object} dtos.JsendErrorResponse "Failure (internal error)"
// @Router /api/images [post]
func (controller *ImageController) PostImages(c *gin.Context) {
	c.Request.ParseMultipartForm(2048)
	var form PostImagesForm
	if err := binding.ShouldBind(&form, c.Request.Form); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(err.(binding.BindingError).FieldErrors))
		return
	}

	if form.ThumbnailUrl == "" {
		form.ThumbnailUrl = form.Url
	}

	if err := controller.imageService.AddImageByURL(form.Url, form.ThumbnailUrl); err != nil {
		if err == utils.FileSizeExceededError {
			c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(map[string]string{
				"url": fmt.Sprintf("Image at url is too large (>%d MB)", config.MAX_IMAGE_FILE_SIZE_MB),
			}))
			return
		} else if err == services.ImageExistsError {
			c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(map[string]string{
				"url": "Image already exists",
			}))
			return
		} else {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, internalErrorJson)
			return
		}
	}

	c.JSON(http.StatusOK, dtos.NewJsendEmptySuccessResponse())
}

type GetImagesQuery struct {
	Offset int `schema:"offset" validate:"min=0"`
	Limit  int `schema:"limit" validate:"min=0"`
}

func ginParamsToMap(params gin.Params) map[string][]string {
	result := make(map[string][]string)
	for _, kv := range params {
		result[kv.Key] = []string{kv.Value}
	}
	return result
}

type ImageIdQuery struct {
	Id int `schema:"id" validate:"min=0"`
}

// @Summary Get image by ID
// @Description Returns an image with the specified ID
// @Tags image
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} dtos.JsendImageResponse "Success"
// @Failure 400 {object} dtos.JsendFailResponse "Failure (bad params)"
// @Failure 404 {object} dtos.JsendFailResponse "Failure (not found)"
// @Failure 500 {object} dtos.JsendErrorResponse "Failure (internal error)"
// @Router /api/images/{id} [get]
func (controller *ImageController) GetImageById(c *gin.Context) {
	var query ImageIdQuery
	if err := binding.ShouldBind(&query, ginParamsToMap(c.Params)); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(err.(binding.BindingError).FieldErrors))
		return
	}
	image, err := controller.imageService.ImageRepo.GetById(query.Id)
	if err == repositories.ImageNotFoundError {
		c.JSON(http.StatusNotFound, dtos.NewJsendFailResponse(map[string]string{
			"id": "No image with such id exists",
		}))
		return
	} else if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	c.JSON(http.StatusOK, dtos.NewJsendImageResponse(*image))
}

// @Summary Delete image by ID
// @Description Deletes an image with the specified ID from the image repository
// @Tags image
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} dtos.JsendEmptySuccessResponse "Successfully deleted image"
// @Failure 400 {object} dtos.JsendFailResponse "Failed to delete image (bad params)"
// @Failure 404 {object} dtos.JsendFailResponse "Failed to delete image (not found)"
// @Failure 500 {object} dtos.JsendErrorResponse "Failed to delete image (internal error)"
// @Router /api/images/{id} [delete]
func (controller *ImageController) DeleteImageById(c *gin.Context) {
	var query ImageIdQuery
	if err := binding.ShouldBind(&query, ginParamsToMap(c.Params)); err != nil {
		c.JSON(http.StatusBadRequest, dtos.NewJsendFailResponse(err.(binding.BindingError).FieldErrors))
		return
	}
	err := controller.imageService.ImageRepo.DeleteById(query.Id)
	if err == repositories.ImageNotFoundError {
		c.JSON(http.StatusNotFound, dtos.NewJsendFailResponse(map[string]string{
			"id": "No image with such id exists",
		}))
		return
	} else if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, internalErrorJson)
		return
	}

	c.JSON(http.StatusOK, dtos.NewJsendEmptySuccessResponse())
}
