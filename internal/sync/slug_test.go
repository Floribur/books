package sync

import "testing"

// TestSlugCollision: two books with title "Dune" → first gets slug "dune",
// second gets slug "dune-1965" (year appended); if years also match → append author surname kebab.
func TestSlugCollision(t *testing.T) {
	// No collision
	slug1 := GenerateSlug("Dune", 1965, "Herbert", map[string]struct{}{})
	if slug1 != "dune" {
		t.Errorf("expected 'dune', got %q", slug1)
	}

	// Collision on base — append year
	slug2 := GenerateSlug("Dune", 1965, "Herbert", map[string]struct{}{"dune": {}})
	if slug2 != "dune-1965" {
		t.Errorf("expected 'dune-1965', got %q", slug2)
	}

	// Collision on base and year — append author surname
	slug3 := GenerateSlug("Dune", 1965, "Herbert", map[string]struct{}{"dune": {}, "dune-1965": {}})
	if slug3 != "dune-herbert" {
		t.Errorf("expected 'dune-herbert', got %q", slug3)
	}
}

// TestGenerateSlug: standard slug generation without collision.
func TestGenerateSlug(t *testing.T) {
	slug := GenerateSlug("The Lord of the Rings", 0, "", map[string]struct{}{})
	if slug != "the-lord-of-the-rings" {
		t.Errorf("expected 'the-lord-of-the-rings', got %q", slug)
	}
}
