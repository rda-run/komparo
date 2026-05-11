package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ReadSnapshot reads a snapshot from a local file or a remote URL.
func ReadSnapshot(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote snapshot: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch remote snapshot: HTTP %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read remote snapshot body: %w", err)
		}
		return data, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read local snapshot: %w", err)
	}
	return data, nil
}
