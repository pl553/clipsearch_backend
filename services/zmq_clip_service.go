package services

type ZmqClipService struct {
	imageEmbeddingEndpoints string
	textEmbeddingEndpoints  string
}

func NewZmqClipService(imageEmbeddingEndpoints string, textEmbeddingEndpoints string) *ZmqClipService {
	return &ZmqClipService{
		imageEmbeddingEndpoints: imageEmbeddingEndpoints,
		textEmbeddingEndpoints:  textEmbeddingEndpoints,
	}
}

func (zcs *ZmqClipService) EncodeImage(imageData []byte) ([]float32, error) {
	conn, err := ConnectToZmqImageEmbeddingDaemon(zcs.imageEmbeddingEndpoints)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.EncodeImage(imageData)
}

func (zcs *ZmqClipService) EncodeText(text string) ([]float32, error) {
	conn, err := ConnectToZmqTextEmbeddingDaemon(zcs.textEmbeddingEndpoints)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.EncodeText(text)
}
