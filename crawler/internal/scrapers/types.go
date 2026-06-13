package scrapers

import (
	"fmt"
	"time"
)

type ScrapedOpportunity struct {
	Title       string
	Type        string
	Country     string
	Deadline    *time.Time
	URL         string
	Source      string
	Description string
	Eligibility string
	Funding     string
}

type ScraperConfig struct {
	Name    string
	BaseURL string
	Type    string
}

var KnownSources = []ScraperConfig{
	{
		Name:    "scholarshipsafrik",
		BaseURL: "https://www.scholarshipsafrik.com",
		Type:    "scholarship",
	},
	{
		Name:    "opportunitiesforafricans",
		BaseURL: "https://www.opportunitiesforafricans.com",
		Type:    "scholarship",
	},
	{
		Name:    "scholarshipstodon",
		BaseURL: "https://scholarshipstodon.com",
		Type:    "scholarship",
	},
}

func FormatDeadline(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func PtrTime(t time.Time) *time.Time {
	return &t
}

func GuessCountry(text string) string {
	countries := map[string]bool{
		"egypt": true, "saudi": true, "uae": true, "qatar": true,
		"kuwait": true, "bahrain": true, "oman": true, "jordan": true,
		"lebanon": true, "morocco": true, "tunisia": true, "algeria": true,
		"iraq": true, "palestine": true, "yemen": true, "libya": true,
		"sudan": true, "somalia": true, "djibouti": true, "comoros": true,
		"mauritania": true, "syria": true,
	}
	lower := text
	for c := range countries {
		if containsIgnoreCase(lower, c) {
			return c
		}
	}
	return ""
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsFold(s, substr))
}

func containsFold(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func DebugLog(source, msg string, args ...any) {
	fmt.Printf("[%s] %s\n", source, fmt.Sprintf(msg, args...))
}
