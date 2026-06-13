package scrapers

import (
	"context"

	"github.com/playwright-community/playwright-go"
)

type Scraper struct {
	Name    string
	BaseURL string
.browser playwright.Browser
	page    playwright.Page
}

func New(name, baseURL string) *Scraper {
	return &Scraper{Name: name, BaseURL: baseURL}
}

func (s *Scraper) Launch(ctx context.Context) error {
	pw, err := playwright.Run()
	if err != nil {
		return err
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return err
	}

	page, err := browser.NewPage()
	if err != nil {
		return err
	}

	s.browser = browser
	s.page = page
	return nil
}

func (s *Scraper) Close() {
	if s.browser != nil {
		s.browser.Close()
	}
}

func (s *Scraper) Navigate(url string) error {
	_, err := s.page.Goto(url)
	return err
}

func (s *Scraper) Content() (string, error) {
	return s.page.Content()
}

func (s *Scraper) Click(selector string) error {
	return s.page.Click(selector)
}

func (s *Scraper) WaitForSelector(selector string) error {
	_, err := s.page.WaitForSelector(selector)
	return err
}
