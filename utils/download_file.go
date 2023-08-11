package utils

import (
	"clipsearch/config"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type FileSizeExceededError struct {
}

func (err FileSizeExceededError) Error() string {
	return fmt.Sprintf("File size exceeded")
}

func (err FileSizeExceededError) Is(target error) bool {
	_, ok := target.(FileSizeExceededError)
	return ok
}

var client = &http.Client{}

// Downloads a file, writing to w
// It checks the content-length header first if it exists
func DownloadFile(w io.Writer, rawUrl string, maxFileSize int) error {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("HEAD", rawUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", url.Hostname())
	req.Header.Set("User-Agent", config.FILE_DOWNLOAD_USERAGENT)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Request wasn't completed successfully, status code %d", resp.StatusCode)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))

	if err == nil && size > maxFileSize {
		return &FileSizeExceededError{}
	}

	req, err = http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", url.Hostname())
	req.Header.Set("User-Agent", config.FILE_DOWNLOAD_USERAGENT)

	resp, err = client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Request wasn't completed successfully, status code %d", resp.StatusCode)
	}

	lr := &LimitedReader{
		Reader:             resp.Body,
		MaxBytesLeftToRead: maxFileSize,
	}

	_, err = io.Copy(w, lr)
	if err != nil {
		return err
	}

	return nil
}

type LimitedReader struct {
	Reader             io.Reader
	MaxBytesLeftToRead int
}

func (lr *LimitedReader) Read(p []byte) (int, error) {
	if len(p) > lr.MaxBytesLeftToRead {
		p = p[0 : lr.MaxBytesLeftToRead+1]
	}
	n, err := lr.Reader.Read(p)
	if err != nil {
		return n, err
	}
	lr.MaxBytesLeftToRead -= n
	if lr.MaxBytesLeftToRead < 0 {
		return n, FileSizeExceededError{}
	}
	return n, nil
}
