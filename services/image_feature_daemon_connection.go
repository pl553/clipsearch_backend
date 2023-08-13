package services

import (
	"encoding/json"
	"fmt"

	"github.com/zeromq/goczmq"
)

type ImageFeatureDaemonConnection struct {
	sock *goczmq.Sock
}

func ConnectToImageFeatureDaemon(endpoints string) (*ImageFeatureDaemonConnection, error) {
	if sock, err := goczmq.NewReq(endpoints); err != nil {
		return nil, err
	} else {
		return &ImageFeatureDaemonConnection{sock: sock}, nil
	}
}

func (conn *ImageFeatureDaemonConnection) SendImageForFeatureExtraction(imageID int, imageData []byte) error {
	err := conn.sock.SendFrame([]byte(fmt.Sprint(imageID)), goczmq.FlagMore)
	if err != nil {
		return err
	}
	err = conn.sock.SendFrame(imageData, goczmq.FlagNone)
	if err != nil {
		return err
	}
	resp, err := conn.sock.RecvMessage()
	if err != nil {
		return err
	}
	var responseJson map[string]any
	if err := json.Unmarshal(resp[0], &responseJson); err != nil {
		return err
	}
	if responseJson["status"] == "error" {
		return fmt.Errorf("Failed to process image: %s", responseJson["message"])
	} else if responseJson["status"] == "fail" {
		return fmt.Errorf("Failed to process image (status: fail)")
	} else if responseJson["status"] == "success" {
		return nil
	} else {
		return fmt.Errorf("Unexpected response format")
	}
}

func (conn *ImageFeatureDaemonConnection) Close() {
	conn.sock.Destroy()
}
