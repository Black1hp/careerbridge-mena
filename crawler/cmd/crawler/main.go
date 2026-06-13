package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/black1hp/careerbridge-mena/crawler/internal/scrapers"
	"github.com/playwright-community/playwright-go"
)

type IngestPayload struct {
	Title       string  `json:"title"`
	Type        string  `json:"type"`
	Country     string  `json:"country"`
	Deadline    *string `json:"deadline"`
	URL         string  `json:"url"`
	Source      string  `json:"source"`
	Description string  `json:"description"`
	Eligibility string  `json:"eligibility"`
	Funding     string  `json:"funding"`
}

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Println("CareerBridge Crawler starting...")

	pw, err := playwright.Run()
	if err != nil {
		fmt.Printf("Playwright error: %v\n", err)
		os.Exit(1)
	}
	defer pw.Stop()

	var allOpps []scrapers.ScrapedOpportunity

	for _, config := range scrapers.KnownSources {
		fmt.Printf("\n=== Scraping: %s ===\n", config.Name)
		scraper := scrapers.NewScholarshipScraper(config)
		opps, err := scraper.Scrape(ctx, pw)
		if err != nil {
			fmt.Printf("Error scraping %s: %v\n", config.Name, err)
			continue
		}
		allOpps = append(allOpps, opps...)
		fmt.Printf("Got %d opportunities from %s\n", len(opps), config.Name)
	}

	if len(allOpps) == 0 {
		fmt.Println("No opportunities scraped. Exiting.")
		return
	}

	fmt.Printf("\n=== Total scraped: %d ===\n", len(allOpps))
	fmt.Println("Ingesting into API...")

	payload := make([]IngestPayload, 0, len(allOpps))
	for _, opp := range allOpps {
		var deadline *string
		if opp.Deadline != nil {
			d := opp.Deadline.Format("2006-01-02")
			deadline = &d
		}
		payload = append(payload, IngestPayload{
			Title:       opp.Title,
			Type:        opp.Type,
			Country:     opp.Country,
			Deadline:    deadline,
			URL:         opp.URL,
			Source:      opp.Source,
			Description: opp.Description,
			Eligibility: opp.Eligibility,
			Funding:     opp.Funding,
		})
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL+"/api/v1/ingest", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Ingest error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Ingest response: %s\n", string(respBody))
	fmt.Println("Done.")
}
