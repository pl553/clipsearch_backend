package models

// swagger:model Image
type Image struct {
	ImageID      int       `json:"id" example:"102"`
	SourceUrl    string    `json:"sourceUrl" example:"http://localhost:8080/example/image.jpg"`
	ThumbnailUrl string    `json:"thumbnailUrl" example:"http://localhost:8080/example/image_thumb.jpg"`
	Sha256       string    `json:"sha256" example:"671797905015849a2e772d7e152ad3289e7d71703b49c8fb607d00265769c1fb"`
	Embedding    []float32 `json:"-"`
}
