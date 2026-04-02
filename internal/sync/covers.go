package sync

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var coverHTTPClient = &http.Client{Timeout: 30 * time.Second}

// ValidateCover checks that data is a valid image >= 5KB and not a 1×1 placeholder.
func ValidateCover(data []byte) error {
	if len(data) < 5*1024 {
		return fmt.Errorf("cover too small: %d bytes (min 5120)", len(data))
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("cover not decodable: %w", err)
	}
	if cfg.Width == 1 && cfg.Height == 1 {
		return fmt.Errorf("cover is 1x1 placeholder")
	}
	return nil
}

// CoverPath returns the local path for a book's cover.
// Uses isbn13 if available, falls back to goodreads_id (D-11, D-12).
func CoverPath(isbn13, goodreadsID string) string {
	if isbn13 != "" {
		return filepath.Join("data", "covers", isbn13+".jpg")
	}
	return filepath.Join("data", "covers", "gr-"+goodreadsID+".jpg")
}

// DownloadCover fetches a cover image URL, validates it, and writes it to destPath.
// Overwrites existing file silently (D-13).
// Returns the local path on success, empty string if validation fails.
func DownloadCover(url, destPath string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("empty cover URL")
	}
	resp, err := coverHTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("download cover: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("rate limited (429)")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cover HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read cover body: %w", err)
	}

	if err := ValidateCover(data); err != nil {
		return "", fmt.Errorf("cover validation: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", fmt.Errorf("create covers dir: %w", err)
	}
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return "", fmt.Errorf("write cover: %w", err)
	}
	log.Printf("cover saved: %s (%d bytes)", destPath, len(data))
	return destPath, nil
}

// TryOpenLibraryCover attempts to download a cover from OpenLibrary by ISBN-13.
// Returns empty string if unavailable or rate-limited.
// Caller is responsible for spacing requests (500ms).
func TryOpenLibraryCover(isbn13, destPath string) string {
	if isbn13 == "" {
		return ""
	}
	url := fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn13)
	path, err := DownloadCover(url, destPath)
	if err != nil {
		log.Printf("openlibrary cover miss for %s: %v", isbn13, err)
		return ""
	}
	return path
}
