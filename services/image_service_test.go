package services

import (
	"clipsearch/repositories"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageService(t *testing.T) {
	t.Run("image add", func(t *testing.T) {
		testImage := struct {
			Path   string
			Sha256 string
		}{
			Path:   "test/test_image.jpg",
			Sha256: "671797905015849a2e772d7e152ad3289e7d71703b49c8fb607d00265769c1fb",
		}

		mockRepo := repositories.NewMockImageRepository()
		imageService := NewImageService(mockRepo)

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			http.ServeFile(rw, req, testImage.Path)
		}))
		defer server.Close()

		err := imageService.AddImageByURL(server.URL)
		if err != nil {
			t.Fatalf(err.Error())
		}
		count, images, err := imageService.GetCountAndImages(0, 1)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if count != 1 {
			t.Fatalf("Image count = %v, want = %v", count, 1)
		}
		if len(images) != 1 {
			t.Fatalf("Returned image count = %v, want = %v", len(images), 1)
		}
		image := images[0]
		if image.Sha256 != testImage.Sha256 {
			t.Fatalf("Returned image sha256 hash = %s, want = %s", image.Sha256, testImage.Sha256)
		}
	})
}
