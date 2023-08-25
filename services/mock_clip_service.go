package services

type MockClipService struct {
}

func NewMockClipService() *MockClipService {
	return &MockClipService{}
}

func (mcs *MockClipService) EncodeImage(imageData []byte) ([]float32, error) {
	return []float32{1, 2, 3}, nil
}

func (mcs *MockClipService) EncodeText(text string) ([]float32, error) {
	return []float32{3, 2, 1}, nil
}
