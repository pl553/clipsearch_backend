package models

type ImageModel struct {
	ImageID   int    `json:"id"`
	SourceUrl string `json:"sourceUrl"`
	Sha256    string `json:"sha256"`
}
