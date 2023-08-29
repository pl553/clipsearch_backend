package dtos

import "clipsearch/models"

// swagger:model JsendImageResponse
type JsendImageResponse struct {
	// Set to "success"
	Status string       `json:"status" example:"success"`
	Data   models.Image `json:"data"`
}

// swagger:model JsendImagesResponse
type JsendImagesResponse struct {
	// Set to "success"
	Status string             `json:"status" example:"success"`
	Data   ImagesResponseData `json:"data"`
}

type ImagesResponseData struct {
	// Total amount of images contained in the repository
	TotalCount int            `json:"totalCount" example:"1234"`
	Images     []models.Image `json:"images"`
}

func NewJsendImageResponse(image models.Image) JsendImageResponse {
	return JsendImageResponse{
		Status: "success",
		Data:   image,
	}
}

func NewJsendImagesResponse(totalCount int, images []models.Image) JsendImagesResponse {
	return JsendImagesResponse{
		Status: "success",
		Data: ImagesResponseData{
			TotalCount: totalCount,
			Images:     images,
		},
	}
}
