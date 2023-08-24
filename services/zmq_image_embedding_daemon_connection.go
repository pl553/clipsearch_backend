package services

import (
	"encoding/json"
	"fmt"

	"github.com/zeromq/goczmq"
)

type ZmqImageEmbeddingDaemonConnection struct {
	sock *goczmq.Sock
}

func ConnectToZmqImageEmbeddingDaemon(endpoints string) (*ZmqImageEmbeddingDaemonConnection, error) {
	if sock, err := goczmq.NewReq(endpoints); err != nil {
		return nil, err
	} else {
		return &ZmqImageEmbeddingDaemonConnection{sock: sock}, nil
	}
}

func (conn *ZmqImageEmbeddingDaemonConnection) EncodeImage(imageData []byte) ([]float32, error) {
	err := conn.sock.SendFrame(imageData, goczmq.FlagNone)
	if err != nil {
		return nil, err
	}

	rawResponse, err := conn.sock.RecvMessage()
	if err != nil {
		return nil, err
	}

	response := struct {
		Status  string
		Message string
		Data    []float32
	}{}

	if err := json.Unmarshal(rawResponse[0], &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		if response.Status == "error" {
			return nil, fmt.Errorf("Failed to process image: %s", response.Message)
		} else {
			return nil, fmt.Errorf("Unexpected response format")
		}
	}

	return response.Data, nil
}

func (conn *ZmqImageEmbeddingDaemonConnection) Close() {
	conn.sock.Destroy()
}
