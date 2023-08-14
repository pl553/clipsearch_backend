package models

type ImageModel struct {
	ImageID      int    `json:"id"`
	SourceUrl    string `json:"sourceUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Sha256       string `json:"sha256"`
}
