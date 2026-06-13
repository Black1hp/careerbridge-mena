package parser

import "time"

type RawOpportunity struct {
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

type Parser interface {
	Parse(html string) ([]RawOpportunity, error)
}
