package d_file

import (
	"fmt"
	"io"
	"net/http"
)

// Bytes downloads the file from the URL and returns its contents as a byte slice.
// Returns nil and an error if the download fails or the URL is empty.
func (f File) Bytes() ([]byte, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("file URL is empty")
	}

	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return bytes, nil
}
