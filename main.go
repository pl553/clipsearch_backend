package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"clipsearch/config"

	"github.com/clevergo/jsend"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgPool *pgxpool.Pool

func getGallery(c *gin.Context) {
	imageUrls := make([]string, 0, 100)
	rows, err := pgPool.Query(context.Background(), `SELECT url FROM images;`)
	if err != nil {
		log.Print("Failed to get image urls from db")
		log.Print(err)
		c.JSON(http.StatusInternalServerError, jsend.NewError("Failed to get images", 199, nil))
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
	router.GET("/api/gallery", getGallery)
	router.Run("localhost:" + port)
}
