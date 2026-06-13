package handlers

import (
	"testing"

	"github.com/black1hp/careerbridge-mena/backend/internal/search"
)

func TestConvertToModels(t *testing.T) {
	opps := []search.ESOpportunity{
		{
			ID:          1,
			Title:       "Test Scholarship",
			Type:        "scholarship",
			Country:     "egypt",
			URL:         "https://example.com",
			Source:      "test",
			Description: "A test scholarship",
		},
	}

	result := convertToModels(opps)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Title != "Test Scholarship" {
		t.Errorf("expected title 'Test Scholarship', got '%s'", result[0].Title)
	}
	if result[0].Type != "scholarship" {
		t.Errorf("expected type 'scholarship', got '%s'", result[0].Type)
	}
}
