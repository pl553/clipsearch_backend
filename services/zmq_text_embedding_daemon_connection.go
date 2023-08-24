package services

import (
	"encoding/json"
	"fmt"

	"github.com/zeromq/goczmq"
)

type ZmqTextEmbeddingDaemonConnection struct {
	sock *goczmq.Sock
}

func ConnectToZmqTextEmbeddingDaemon(endpoints string) (*ZmqTextEmbeddingDaemonConnection, error) {
	if sock, err := goczmq.NewReq(endpoints); err != nil {
		return nil, err
	} else {
		return &ZmqTextEmbeddingDaemonConnection{sock: sock}, nil
	}
}

func (conn *ZmqTextEmbeddingDaemonConnection) EncodeText(text string) ([]float32, error) {
	err := conn.sock.SendFrame([]byte(text), goczmq.FlagNone)
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
			return nil, fmt.Errorf("Failed to process text: %s", response.Message)
		} else {
			return nil, fmt.Errorf("Unexpected response format")
		}
	}

	return response.Data, nil
}

func (conn *ZmqTextEmbeddingDaemonConnection) Close() {
	conn.sock.Destroy()
}
