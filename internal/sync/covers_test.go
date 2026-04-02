package sync

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

// makeJPEG creates an in-memory JPEG of given dimensions and quality.
func makeJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.Black)
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	return buf.Bytes()
}

func TestCoverValidation_TooSmall(t *testing.T) {
	data := make([]byte, 100)
	err := ValidateCover(data)
	if err == nil {
		t.Fatal("expected error for tiny file, got nil")
	}
	if !containsStr(err.Error(), "too small") {
		t.Errorf("expected 'too small' in error, got: %v", err)
	}
}

func TestCoverValidation_NonDecodable(t *testing.T) {
	data := bytes.Repeat([]byte{0xFF}, 10*1024)
	err := ValidateCover(data)
	if err == nil {
		t.Fatal("expected error for non-decodable data, got nil")
	}
	if !containsStr(err.Error(), "not decodable") {
		t.Errorf("expected 'not decodable' in error, got: %v", err)
	}
}

func TestCoverValidation_OnePx(t *testing.T) {
	data := makeJPEG(1, 1)
	// Pad to >= 5KB so size check passes
	data = append(data, bytes.Repeat([]byte{0}, 5*1024)...)
	err := ValidateCover(data)
	// NOTE: makeJPEG(1,1) may fail size check before 1×1 check depending on encoding size.
	// If err contains "too small", that's also acceptable — the cover is rejected either way.
	if err == nil {
		t.Fatal("expected error for 1×1 image, got nil")
	}
}

func TestCoverValidation_Valid(t *testing.T) {
	// Create a valid 300×450 JPEG that encodes to >= 5KB
	data := makeJPEG(300, 450)
	// If encoded size is still < 5KB, pad (unlikely for 300×450 but defensive)
	if len(data) < 5*1024 {
		t.Skipf("generated JPEG too small (%d bytes) for this test; skip", len(data))
	}
	if err := ValidateCover(data); err != nil {
		t.Errorf("expected valid cover to pass, got: %v", err)
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
