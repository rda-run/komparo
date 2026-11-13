package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rda-run/komparo/pkg/schema"
)

// The default http.Client has no timeout and would hang a CI job forever.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// ReadSnapshot loads and parses a snapshot from a local file or a remote URL.
func ReadSnapshot(path string) (*schema.SchemaSnapshot, error) {
	data, err := readBytes(path)
	if err != nil {
		return nil, err
	}

	var snapshot schema.SchemaSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot JSON: %w", err)
	}
	snapshot.Normalize()
	return &snapshot, nil
}

func readBytes(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := httpClient.Get(path)
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
