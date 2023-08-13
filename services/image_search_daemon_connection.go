package services

import (
	"encoding/json"
	"fmt"

	"github.com/zeromq/goczmq"
)

type ImageSearchResult struct {
	ImageID int     `json:"id"`
	Score   float32 `json:"score"`
}

type ImageSearchDaemonConnection struct {
	sock *goczmq.Sock
}

func ConnectToImageSearchDaemon(endpoints string) (*ImageSearchDaemonConnection, error) {
	if sock, err := goczmq.NewReq(endpoints); err != nil {
		return nil, err
	} else {
		return &ImageSearchDaemonConnection{sock: sock}, nil
	}
}

func (conn *ImageSearchDaemonConnection) SearchRequest(query string) (int, []ImageSearchResult, error) {
	err := conn.sock.SendFrame([]byte(query), goczmq.FlagNone)
	if err != nil {
		return 0, nil, err
	}
	resp, err := conn.sock.RecvMessage()
	if err != nil {
		return 0, nil, err
	}
	var responseJson map[string]any
	if err := json.Unmarshal(resp[0], &responseJson); err != nil {
		return 0, nil, err
	}

	if responseJson["status"] != "success" {
		if responseJson["status"] == "error" {
			return 0, nil, fmt.Errorf("Failed to process image: %s", responseJson["message"])
		} else if responseJson["status"] == "fail" {
			return 0, nil, fmt.Errorf("Failed to process image (status: fail)")
		} else {
			return 0, nil, fmt.Errorf("Unexpected response format")
		}
	}

	searchResults := struct {
		Status string
		Data   struct {
			TotalCount int
			Images     []ImageSearchResult
		}
	}{}

	if err := json.Unmarshal(resp[0], &searchResults); err != nil {
		return 0, nil, err
	}

	return searchResults.Data.TotalCount, searchResults.Data.Images, nil
}

func (conn *ImageSearchDaemonConnection) Close() {
	conn.sock.Destroy()
}
