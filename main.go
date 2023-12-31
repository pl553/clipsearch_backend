package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"clipsearch/config"
	"clipsearch/controllers"
	"clipsearch/repositories"
	"clipsearch/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func seedRepository(imageRepository repositories.ImageRepository, imageService *services.ImageService) {
	count, err := imageRepository.Count()
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		seedImageURLs := []string{
			"http://localhost/static/images/1.gif",
			"http://localhost/static/images/2.jpg",
			"http://localhost/static/images/3.jpg",
		}

		for _, imageURL := range seedImageURLs {
			if err := imageService.AddImageByURL(imageURL, ""); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func setupRouter(imageController *controllers.ImageController) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/api/images", imageController.GetImages)
	router.POST("/api/images", imageController.PostImages)
	router.GET("/api/images/:id", imageController.GetImageById)
	router.DELETE("/api/images/:id", imageController.DeleteImageById)
	router.GET("/api/images/search", imageController.GetSearchImages)
	return router
}

// @title CLIP search API
// @version         1.0
func main() {
	port := os.Getenv(config.PORT_ENVAR)
	if port == "" {
		port = config.DEFAULT_PORT
	}
	zmq_text_port := os.Getenv(config.ZMQ_TEXT_EMBEDDING_DAEMON_PORT_ENVAR)
	if zmq_text_port == "" {
		zmq_text_port = config.ZMQ_TEXT_EMBEDDING_DAEMON_DEFAULT_PORT
	}
	zmq_image_port := os.Getenv(config.ZMQ_IMAGE_EMBEDDING_DAEMON_PORT_ENVAR)
	if zmq_image_port == "" {
		zmq_image_port = config.ZMQ_IMAGE_EMBEDDING_DAEMON_DEFAULT_PORT
	}
	dbConnString := os.Getenv(config.PG_DATABASE_CONNECTION_URL_ENVAR)
	if dbConnString == "" {
		log.Fatal(fmt.Sprintf("Please define the %v envar", config.PG_DATABASE_CONNECTION_URL_ENVAR))
	}
	pgPool, err := pgxpool.New(context.Background(), dbConnString)
	if err != nil {
		log.Print("Failed to connect to db!")
		log.Fatal(err)
	}
	defer pgPool.Close()

	clipService := services.NewZmqClipService("tcp://localhost:"+zmq_image_port, "tcp://localhost:"+zmq_text_port)

	imageRepository := repositories.NewPgImageRepository(pgPool)
	imageService := services.NewImageService(imageRepository, clipService)
	imageController := controllers.NewImageController(imageService)

	//seedRepository(imageRepository, imageService)

	router := setupRouter(imageController)
	router.Run(":" + port)
}
