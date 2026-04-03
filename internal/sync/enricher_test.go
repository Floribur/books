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
			want:            true,
		},
		{
			name:            "case insensitive match passes",
			inputTitle:      "the hobbit",
			inputAuthor:     "tolkien",
			returnedTitle:   "The Hobbit",
			returnedAuthors: []string{"J.R.R. Tolkien"},
			want:            true,
		},
		// Real-world cases from Google Books enricher logs
		{
			name:            "goodreads full title vs abbreviated google title passes",
			inputTitle:      "Modeling Mindsets: The Many Cultures Of Learning From Data",
			inputAuthor:     "",
			returnedTitle:   "MODELING MINDSETS",
			returnedAuthors: nil,
			want:            true,
		},
		{
			name:            "goodreads subtitle stripped vs short google title passes",
			inputTitle:      "Halt die Klappe, Kopf! Ein Selfcare-Buch für Tage, an denen Schokokuchen nicht reicht",
			inputAuthor:     "",
			returnedTitle:   "Halt die Klappe, Kopf!",
			returnedAuthors: nil,
			want:            true,
		},
		{
			name:            "sapiens subtitle vs short google title passes",
			inputTitle:      "Sapiens: A Brief History of Humankind",
			inputAuthor:     "",
			returnedTitle:   "Sapiens",
			returnedAuthors: nil,
			want:            true,
		},
		{
			name:            "series book vs series title fails",
			inputTitle:      "Mockingjay (The Hunger Games, #3)",
			inputAuthor:     "",
			returnedTitle:   "The Hunger Games",
			returnedAuthors: nil,
			want:            false,
		},
		{
			name:            "completely wrong result fails",
			inputTitle:      "The Hunger Games (Hunger Games, #1)",
			inputAuthor:     "",
			returnedTitle:   "Beyond the Screen",
			returnedAuthors: nil,
			want:            false,
		},
		{
			name:            "hyphen vs space in title passes",
			inputTitle:      "Game Over - der Fall der Credit Suisse: Das Buch zum gleichnamigen Doku-Film (German Edition)",
			inputAuthor:     "",
			returnedTitle:   "GAME-OVER - der Fall der Credit Suisse",
			returnedAuthors: nil,
			want:            true,
		},
		{
			name:            "google uses series name as title passes",
			inputTitle:      "The 11:59 Bomber (NYPD Red, #8)",
			inputAuthor:     "James Patterson",
			returnedTitle:   "NYPD Red 8",
			returnedAuthors: []string{"James Patterson"},
			want:            true,
		},
		{
			name:            "totally wrong result for series book fails",
			inputTitle:      "The Strength of the Few (Hierarchy, #2)",
			inputAuthor:     "",
			returnedTitle:   "The American Gas Light Journal",
			returnedAuthors: nil,
			want:            false,
		},
		{
			name:            "german translation returned instead of original fails",
			inputTitle:      "Deadly Heat (Nikki Heat, #5)",
			inputAuthor:     "Richard Castle",
			returnedTitle:   "Castle 5: Deadly Heat - Tödliche Hitze",
			returnedAuthors: []string{"Richard Castle"},
			want:            false,
		},
		{
			name:            "different book in same series fails",
			inputTitle:      "An Assassin's Creed Series. Last Descendants. Aufstand in New York (German Edition)",
			inputAuthor:     "",
			returnedTitle:   "An Assassin's Creed Series. Last Descendants. Das Grab des Khan",
			returnedAuthors: nil,
			want:            false,
		},
		{
			name:            "double space in author name normalizes correctly",
			inputTitle:      "NYPD Red 2 (NYPD Red, #2)",
			inputAuthor:     "James  Patterson", // Goodreads RSS often has double spaces
			returnedTitle:   "NYPD Red 2",
			returnedAuthors: []string{"James Patterson", "Marshall Karp"},
			want:            true,
		},
		{
			name:            "single goodreads author matches when google has co-author",
			inputTitle:      "NYPD Red 2",
			inputAuthor:     "James Patterson",
			returnedTitle:   "NYPD Red 2",
			returnedAuthors: []string{"James Patterson", "Marshall Karp"},
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
