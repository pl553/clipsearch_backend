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

	"github.com/clevergo/jsend"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgPool *pgxpool.Pool

type PostGalleryImagesForm struct {
	Url string `schema:"url,required" validate:"url"`
}

func postGalleryImages(c *gin.Context) {
	c.Request.ParseMultipartForm(2048)
	var form PostGalleryImagesForm
	if err := binding.ShouldBind(&form, c.Request.Form); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err))
		return
	}
	rawUrl := c.Request.FormValue("url")
	if rawUrl == "" {
		c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{"url": "Image url must be non-empty"}))
		return
	}
	invalidUrlMessage := jsend.NewFail(gin.H{"url": "Invalid url"})
	url, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, invalidUrlMessage)
		return
	}
	if len(url.Query()) != 0 {
		c.JSON(http.StatusBadRequest, jsend.NewFail(gin.H{"url": "Url must not have any query parameters in it"}))
		return
	}

	// TODO add more validation (check that its actually an image, check that the filesize isnt too large, ...)
	query := `INSERT INTO images (url) VALUES ($1);`
	rows, err := pgPool.Query(context.Background(), query, rawUrl)
	rows.Close()
	if err != nil {
		log.Printf("POST /api/gallery/images: failed to execute query %v", query)
		log.Print(err)
		c.JSON(http.StatusInternalServerError, jsend.NewError("Internal error", 500, nil))
		return
	}

	c.JSON(http.StatusOK, jsend.New(nil))
}

type GetGalleryImagesQuery struct {
	Offset int `schema:"offset" validate:"min=0"`
	Limit  int `schema:"limit" validate:"min=0"`
}

func getGalleryImages(c *gin.Context) {
	var query GetGalleryImagesQuery

	if err := binding.ShouldBind(&query, c.Request.URL.Query()); err != nil {
		c.JSON(http.StatusBadRequest, jsend.NewFail(err))
		return
	}

	imageUrls := make([]string, 0, 100)

	sqlCountQuery := `SELECT COUNT(*) FROM images`
	row := pgPool.QueryRow(context.Background(), sqlCountQuery)
	var count int
	if err := row.Scan(&count); err != nil {
		c.JSON(http.StatusInternalServerError, jsend.NewError("Failed to get images", 500, nil))
		return
	}

	sqlQuery := `SELECT image_id, url FROM images ORDER BY image_id LIMIT $1 OFFSET $2;`
	rows, err := pgPool.Query(context.Background(), sqlQuery, query.Limit, query.Offset)
	defer rows.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.NewError("Failed to get images", 500, nil))
		return
	}
	for rows.Next() {
		var url string
		rows.Scan(nil, &url)
		imageUrls = append(imageUrls, url)
	}
	rows.Close()
	c.JSON(http.StatusOK, jsend.New(gin.H{
		"image_count": count,
		"image_urls":  imageUrls,
	}))
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
	var err error
	pgPool, err = pgxpool.New(context.Background(), dbConnString)
	if err != nil {
		log.Print("Failed to connect to db!")
		log.Fatal(err)
	}
	defer pgPool.Close()

	rows, err := pgPool.Query(context.Background(), "SELECT * FROM images;")
	if err != nil {
		log.Fatal("Failed to execute query")
	}

	if !rows.Next() {
		_, err := pgPool.Query(context.Background(), `INSERT INTO images (url)
		   VALUES
		     ('/static/images/1.gif'),
		     ('/static/images/2.jpg'),
		     ('/static/images/3.jpg');`)
		if err != nil {
			log.Fatal("Failed to seed db")
		}
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/api/gallery/images", getGalleryImages)
	router.POST("/api/gallery/images", postGalleryImages)
	router.Run("localhost:" + port)
}
