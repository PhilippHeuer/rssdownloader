package feed

import (
	"fmt"
	"io"
	"os"
)

const maxFileSize = 500 * 1024 * 1024

func DownloadFile(url, filename string, timeoutSeconds *int64) error {
	resp, err := GetClient(timeoutSeconds).Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.ContentLength > maxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", resp.ContentLength, maxFileSize)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	limitedReader := io.LimitReader(resp.Body, maxFileSize)
	_, err = io.Copy(file, limitedReader)
	return err
}
