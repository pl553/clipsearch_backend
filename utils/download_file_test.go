package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	t.Run("small file", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("OK"))
		}))
		defer server.Close()

		var buf bytes.Buffer
		err := DownloadFile(&buf, server.URL, 1)
		if err != FileSizeExceededError {
			t.Fatalf("Expected download to fail with error FileSizeExceeded")
		}
		buf.Reset()
		err = DownloadFile(&buf, server.URL, 2)
		if err != nil {
			t.Fatalf(err.Error())
		}
		content := buf.Bytes()
		if err != nil {
			t.Fatalf(err.Error())
		}
		if string(content) != "OK" {
			t.Fatalf(`Downloaded file contents = %s, want = "OK"`, string(content))
		}
	})

	t.Run("big file", func(t *testing.T) {
		serverFileContent := bytes.Repeat([]byte("1"), 500*1024)
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write(serverFileContent)
		}))
		var buf bytes.Buffer
		err := DownloadFile(&buf, server.URL, 500*1024)
		if err != nil {
			t.Fatalf(err.Error())
		}

		content := buf.Bytes()
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !bytes.Equal(content, serverFileContent) {
			t.Fatalf("Downloaded file contents dont match")
		}
	})

	t.Run("content length", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Content-Length", "16") // it should check content-length first
			rw.Write([]byte("OK"))
		}))
		var buf bytes.Buffer
		err := DownloadFile(&buf, server.URL, 8)
		if err != FileSizeExceededError {
			t.Fatalf("Expected download to fail with error FileSizeExceeded")
		}
	})
}
