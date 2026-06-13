package scrapers

import "testing"

func TestGuessCountry(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Scholarship in Egypt for students", "egypt"},
		{"Saudi Arabia internship", "saudi"},
		{"UAE competition 2026", "uae"},
		{"Random text with no country", ""},
		{"Jordan scholarship program", "jordan"},
		{"Lebanese university grant", "lebanon"},
	}

	for _, tt := range tests {
		got := GuessCountry(tt.input)
		if got != tt.want {
			t.Errorf("GuessCountry(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGuessType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Full scholarship for masters", "scholarship"},
		{"Summer internship at Google", "internship"},
		{"Hackathon competition 2026", "competition"},
		{"Grant for research", "scholarship"},
	}

	for _, tt := range tests {
		got := guessType(tt.input)
		if got != tt.want {
			t.Errorf("guessType(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExtractDeadline(t *testing.T) {
	deadline := extractDeadline("Deadline: 15 Jan 2026")
	if deadline == nil {
		t.Error("expected deadline, got nil")
		return
	}
	if deadline.Day() != 15 || deadline.Month() != 1 {
		t.Errorf("expected 15 Jan, got %v", deadline)
	}

	none := extractDeadline("no date here")
	if none != nil {
		t.Errorf("expected nil, got %v", none)
	}
}

func TestExtractFunding(t *testing.T) {
	funding := extractFunding("This is a fully funded scholarship with stipend")
	if funding == "" {
		t.Error("expected funding info, got empty")
	}
}

func TestExtractEligibility(t *testing.T) {
	elig := extractEligibility("Eligibility: Must be under 30 years old")
	if elig == "" {
		t.Error("expected eligibility info, got empty")
	}
}
