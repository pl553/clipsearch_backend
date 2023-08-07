package models

type ImageModel struct {
	Id        int    `db:"id"`
	SourceUrl string `db:"source_url"`
}
