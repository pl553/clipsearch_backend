package controllers

import (
	"clipsearch/models"
	"clipsearch/repositories"
	"clipsearch/services"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testImage = struct {
	Path   string
	Sha256 string
}{
	Path:   "../test/test_image.jpg",
	Sha256: "671797905015849a2e772d7e152ad3289e7d71703b49c8fb607d00265769c1fb",
}

func TestImageController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testImageServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, testImage.Path)
	}))
	defer testImageServer.Close()

	t.Run("PostImages", func(t *testing.T) {
		mockRepo := repositories.NewMockImageRepository()
		imageService := services.NewImageService(mockRepo)
		controller := NewImageController(imageService)

		router := gin.Default()
		router.POST("/api/images", controller.PostImages)

		t.Run("should return 400 if url is not provided", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/images", strings.NewReader("url="))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 200 if theres an image at url", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/images", strings.NewReader("url="+testImageServer.URL))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)
			var result map[string]any
			err := json.Unmarshal(resp.Body.Bytes(), &result)
			assert.Equal(t, nil, err)
			assert.Equal(t, "success", result["status"])
		})
	})

	t.Run("GetImages", func(t *testing.T) {
		t.Run("empty repo", func(t *testing.T) {
			mockRepo := repositories.NewMockImageRepository()
			imageService := services.NewImageService(mockRepo)
			controller := NewImageController(imageService)

			router := gin.Default()
			router.GET("/api/images", controller.GetImages)

			req, _ := http.NewRequest(http.MethodGet, "/api/images", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)

			result := struct {
				Status string
				Data   struct {
					ImageCount int
				}
			}{}
			err := json.Unmarshal(resp.Body.Bytes(), &result)
			assert.Equal(t, nil, err)

			assert.Equal(t, "success", result.Status)
			assert.Equal(t, 0, result.Data.ImageCount)
		})

		t.Run("repo with images", func(t *testing.T) {
			mockRepo := repositories.NewMockImageRepository()
			imageService := services.NewImageService(mockRepo)

			imageService.AddImageByURL(testImageServer.URL)
			imageService.AddImageByURL(testImageServer.URL)
			imageService.AddImageByURL(testImageServer.URL)

			controller := NewImageController(imageService)

			router := gin.Default()
			router.GET("/api/images", controller.GetImages)

			req, _ := http.NewRequest(http.MethodGet, "/api/images?offset=1&limit=2", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)

			result := struct {
				Status string
				Data   struct {
					ImageCount int
					Images     []models.ImageModel
				}
			}{}
			err := json.Unmarshal(resp.Body.Bytes(), &result)
			assert.Equal(t, nil, err)

			assert.Equal(t, "success", result.Status)
			assert.Equal(t, 3, result.Data.ImageCount)
			assert.Equal(t, 2, len(result.Data.Images))
			assert.Equal(t, testImage.Sha256, result.Data.Images[1].Sha256)
		})
	})

	t.Run("GetImageById", func(t *testing.T) {
		t.Run("should return 404 if id is not provided", func(t *testing.T) {
			mockRepo := repositories.NewMockImageRepository()
			imageService := services.NewImageService(mockRepo)
			controller := NewImageController(imageService)

			router := gin.Default()
			router.GET("/api/images/:id", controller.GetImageById)

			req, _ := http.NewRequest(http.MethodGet, "/api/images/", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		t.Run("should return data if id is valid", func(t *testing.T) {
			mockRepo := repositories.NewMockImageRepository()
			imageService := services.NewImageService(mockRepo)

			imageService.AddImageByURL(testImageServer.URL)

			controller := NewImageController(imageService)

			router := gin.Default()
			router.GET("/api/images/:id", controller.GetImageById)

			req, _ := http.NewRequest(http.MethodGet, "/api/images/1", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)
			result := struct {
				Status string
				Data   models.ImageModel
			}{}
			err := json.Unmarshal(resp.Body.Bytes(), &result)
			assert.Equal(t, nil, err)
			log.Print(result)
			assert.Equal(t, "success", result.Status)
			assert.Equal(t, testImage.Sha256, result.Data.Sha256)
			assert.Equal(t, 1, result.Data.ImageID)
			assert.Equal(t, testImageServer.URL, result.Data.SourceUrl)
		})
	})
}
