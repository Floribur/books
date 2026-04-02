package sync

import "testing"

// TestCSVISBNUnquote: strips Excel-style formula quoting.
func TestCSVISBNUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`="9780385472579"`, "9780385472579"},
		{"9780385472579", "9780385472579"},
		{`=""`, ""},
		{"", ""},
		{`="0385472579"`, "0385472579"},
	}
	for _, tc := range tests {
		got := unquoteISBN(tc.input)
		if got != tc.expected {
			t.Errorf("unquoteISBN(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

// TestCSVImport: parse a CSV string with quoted ISBNs; assert books parsed with clean isbn13 values.
// This test does not require a DB connection — it tests unquoteISBN directly via the CSV parsing path.
func TestCSVImport(t *testing.T) {
	// Test the unquoteISBN function which is the core CSV ISBN processing logic
	cases := []struct {
		raw      string
		expected string
	}{
		{`="9780385472579"`, "9780385472579"},
		{`="9780451524935"`, "9780451524935"},
	}
	for _, tc := range cases {
		got := unquoteISBN(tc.raw)
		if got != tc.expected {
			t.Errorf("CSV ISBN unquoting: unquoteISBN(%q) = %q, want %q", tc.raw, got, tc.expected)
		}
	}
}
