package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"clipsearch/config"

	"github.com/clevergo/jsend"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgPool *pgxpool.Pool

func postGalleryImages(c *gin.Context) {
	c.Request.ParseMultipartForm(2048)
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
	_, err = pgPool.Query(context.Background(), query, rawUrl)
	if err != nil {
		log.Printf("POST /api/gallery/images: failed to execute query %v", query)
		log.Print(err)
		c.JSON(http.StatusInternalServerError, jsend.NewError("Internal error", 500, nil))
		return
	}
	c.JSON(http.StatusOK, jsend.New(nil))
}

func getGalleryImages(c *gin.Context) {
	imageUrls := make([]string, 0, 100)
	rows, err := pgPool.Query(context.Background(), `SELECT url FROM images;`)
	if err != nil {
		log.Print("Failed to get image urls from db")
		log.Print(err)
		c.JSON(http.StatusInternalServerError, jsend.NewError("Failed to get images", 500, nil))
		return
	}
	for rows.Next() {
		var url string
		rows.Scan(&url)
		imageUrls = append(imageUrls, url)
	}
	c.JSON(http.StatusOK, jsend.New(gin.H{"image_urls": imageUrls}))
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
		     ('http://localhost/static/images/1.gif'),
		     ('http://localhost/static/images/2.jpg'),
		     ('http://localhost/static/images/3.jpg');`)
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
