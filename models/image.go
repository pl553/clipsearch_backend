package models

type ImageModel struct {
	Id        int    `db:"id" json:"id"`
	SourceUrl string `db:"source_url" json:"source_url"`
}
