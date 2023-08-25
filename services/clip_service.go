package services

type ClipService interface {
	EncodeImage(imageData []byte) ([]float32, error)
	EncodeText(text string) ([]float32, error)
}
