package main

import (
	"os"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/clevergo/jsend"

	"clipsearch/config"
)

func getGallery(c *gin.Context) {
	c.JSON(http.StatusOK, jsend.New(gin.H{"image_urls": config.ImageUrls}))
}

func main() {
	port := os.Getenv(config.PortEnvar)
	if port == "" {
		port = config.DefaultPort
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/api/gallery", getGallery)
	router.Run("localhost:" + port)
}
