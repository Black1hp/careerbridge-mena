package scrapers

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type ScholarshipScraper struct {
	Config ScraperConfig
}

func NewScholarshipScraper(config ScraperConfig) *ScholarshipScraper {
	return &ScholarshipScraper{Config: config}
}

func (s *ScholarshipScraper) Scrape(ctx context.Context, pw *playwright.Playwright) ([]ScrapedOpportunity, error) {
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: boolPtr(true),
	})
	if err != nil {
		return nil, fmt.Errorf("launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("new page: %w", err)
	}
	defer page.Close()

	DebugLog(s.Config.Name, "Navigating to %s", s.Config.BaseURL)
	_, err = page.Goto(s.Config.BaseURL, playwright.PageGotoOptions{
		Timeout:   float64Ptr(30000),
		WaitUntil: playwright.WaitUntilNetworkIdle,
	})
	if err != nil {
		return nil, fmt.Errorf("goto: %w", err)
	}

	var opportunities []ScrapedOpportunity

	articleSelector := "article, .post, .entry, .scholarship-item, .opportunity, [class*='scholar'], [class*='opport'], [class*='list-item']"
	articles, err := page.QuerySelectorAll(articleSelector)
	if err != nil || len(articles) == 0 {
		DebugLog(s.Config.Name, "No articles found with primary selector, trying fallback")
		articles, _ = page.QuerySelectorAll("a[href*='scholar'], a[href*='opport'], a[href*='intern'], .card, .item")
	}

	DebugLog(s.Config.Name, "Found %d articles", len(articles))

	for i, article := range articles {
		if i >= 50 {
			break
		}

		titleEl, _ := article.QuerySelector("h1, h2, h3, h4, .title, .entry-title, a")
		if titleEl == nil {
			continue
		}
		title, _ := titleEl.TextContent()
		title = strings.TrimSpace(title)
		if title == "" {
			continue
		}

		linkEl, _ := article.QuerySelector("a[href]")
		var link string
		if linkEl != nil {
			link, _ = linkEl.GetAttribute("href")
		}
		if link == "" {
			link = s.Config.BaseURL
		}
		if !strings.HasPrefix(link, "http") {
			link = s.Config.BaseURL + link
		}

		descEl, _ := article.QuerySelector("p, .excerpt, .summary, .description, .entry-content")
		var description string
		if descEl != nil {
			description, _ = descEl.TextContent()
			description = strings.TrimSpace(description)
		}

		deadline := extractDeadline(title + " " + description)
		country := GuessCountry(title + " " + description)

		oppType := s.Config.Type
		if oppType == "" {
			oppType = guessType(title + " " + description)
		}

		opportunities = append(opportunities, ScrapedOpportunity{
			Title:       title,
			Type:        oppType,
			Country:     country,
			Deadline:    deadline,
			URL:         link,
			Source:      s.Config.Name,
			Description: description,
			Eligibility: extractEligibility(description),
			Funding:     extractFunding(description),
		})
	}

	DebugLog(s.Config.Name, "Scraped %d opportunities", len(opportunities))
	return opportunities, nil
}

func extractDeadline(text string) *time.Time {
	patterns := []string{
		`deadline[:\s]*(\d{1,2}[\s/\-\.](?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*[\s/\-\.]\d{4})`,
		`deadline[:\s]*(\d{1,2}[\s/\-\.]\d{1,2}[\s/\-\.]\d{4})`,
		`(\d{1,2}[\s/\-\.](?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*[\s/\-\.]\d{4})`,
		`(\d{4}[\s/\-\.]\d{1,2}[\s/\-\.]\d{1,2})`,
	}

	dateFormats := []string{
		"02 Jan 2006", "02 January 2006",
		"02/01/2006", "02-01-2006", "02.01.2006",
		"2006-01-02", "2006/01/02",
		"Jan 02, 2006", "January 02, 2006",
		"02 Jan, 2006", "02 January, 2006",
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			dateStr := strings.TrimSpace(matches[1])
			for _, format := range dateFormats {
				if t, err := time.Parse(format, dateStr); err == nil {
					return &t
				}
			}
		}
	}

	return nil
}

func guessType(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "scholarship") || strings.Contains(lower, "grant") || strings.Contains(lower, "fellowship") {
		return "scholarship"
	}
	if strings.Contains(lower, "internship") || strings.Contains(lower, "intern") {
		return "internship"
	}
	if strings.Contains(lower, "competition") || strings.Contains(lower, "contest") || strings.Contains(lower, "hackathon") {
		return "competition"
	}
	return "scholarship"
}

func extractEligibility(text string) string {
	lower := strings.ToLower(text)
	keywords := []string{"eligible", "eligibility", "requirements", "who can apply", "criteria", "qualification"}
	for _, kw := range keywords {
		idx := strings.Index(lower, kw)
		if idx >= 0 {
			start := idx
			if start > 20 {
				start -= 20
			}
			end := idx + 200
			if end > len(text) {
				end = len(text)
			}
			return strings.TrimSpace(text[start:end])
		}
	}
	return ""
}

func extractFunding(text string) string {
	lower := strings.ToLower(text)
	keywords := []string{"funding", "coverage", "stipend", "tuition", "financial", "fully funded", "partially funded"}
	for _, kw := range keywords {
		idx := strings.Index(lower, kw)
		if idx >= 0 {
			start := idx
			if start > 10 {
				start -= 10
			}
			end := idx + 150
			if end > len(text) {
				end = len(text)
			}
			return strings.TrimSpace(text[start:end])
		}
	}
	return ""
}

func boolPtr(b bool) *bool { return &b }
func float64Ptr(f float64) *float64 { return &f }
