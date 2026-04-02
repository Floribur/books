package sync

import "testing"

func TestEnrichmentConfidenceGate(t *testing.T) {
	tests := []struct {
		name            string
		inputTitle      string
		inputAuthor     string
		returnedTitle   string
		returnedAuthors []string
		want            bool
	}{
		{
			name:            "exact match passes",
			inputTitle:      "Dune",
			inputAuthor:     "Frank Herbert",
			returnedTitle:   "Dune",
			returnedAuthors: []string{"Frank Herbert"},
			want:            true,
		},
		{
			name:            "completely wrong book fails",
			inputTitle:      "Dune",
			inputAuthor:     "Frank Herbert",
			returnedTitle:   "Foundation",
			returnedAuthors: []string{"Isaac Asimov"},
			want:            false,
		},
		{
			name:            "title matches but wrong author fails",
			inputTitle:      "Dune",
			inputAuthor:     "Frank Herbert",
			returnedTitle:   "Dune",
			returnedAuthors: []string{"Isaac Asimov"},
			want:            false,
		},
		{
			name:            "returned title contains input title passes",
			inputTitle:      "Dune",
			inputAuthor:     "Herbert",
			returnedTitle:   "Dune Messiah",
			returnedAuthors: []string{"Frank Herbert"},
			want:            true, // "dune messiah" contains "dune"
		},
		{
			name:            "case insensitive match passes",
			inputTitle:      "the hobbit",
			inputAuthor:     "tolkien",
			returnedTitle:   "The Hobbit",
			returnedAuthors: []string{"J.R.R. Tolkien"},
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := confidenceGate(tt.inputTitle, tt.inputAuthor, tt.returnedTitle, tt.returnedAuthors)
			if got != tt.want {
				t.Errorf("confidenceGate(%q, %q, %q, %v) = %v, want %v",
					tt.inputTitle, tt.inputAuthor, tt.returnedTitle, tt.returnedAuthors, got, tt.want)
			}
		})
	}
}
