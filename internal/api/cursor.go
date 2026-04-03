package api

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// encodeCursor encodes read_at + id into an opaque base64 URL-safe cursor token.
// Format before encoding: "<RFC3339Nano>_<id>"
func encodeCursor(readAt time.Time, id int64) string {
	raw := fmt.Sprintf("%s_%d", readAt.UTC().Format(time.RFC3339Nano), id)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// decodeCursor decodes a cursor token produced by encodeCursor.
// Returns an error if the token is malformed or invalid base64.
func decodeCursor(s string) (time.Time, int64, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid cursor: %w", err)
	}
	parts := strings.SplitN(string(b), "_", 2)
	if len(parts) != 2 {
		return time.Time{}, 0, fmt.Errorf("invalid cursor format")
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid cursor timestamp: %w", err)
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("invalid cursor id: %w", err)
	}
	return t, id, nil
}

// Export for testing (test package imports via api.EncodeCursor / api.DecodeCursor).
var EncodeCursor = encodeCursor
var DecodeCursor = decodeCursor
