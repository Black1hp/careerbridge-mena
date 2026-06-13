package models

import "time"

type OpportunityType string

const (
	TypeScholarship  OpportunityType = "scholarship"
	TypeInternship   OpportunityType = "internship"
	TypeCompetition  OpportunityType = "competition"
)

type Opportunity struct {
	ID          int             `json:"id" db:"id"`
	Title       string          `json:"title" db:"title"`
	Type        OpportunityType `json:"type" db:"type"`
	Country     string          `json:"country" db:"country"`
	Deadline    *time.Time      `json:"deadline" db:"deadline"`
	URL         string          `json:"url" db:"url"`
	Source      string          `json:"source" db:"source"`
	Description string          `json:"description" db:"description"`
	Eligibility string          `json:"eligibility" db:"eligibility"`
	Funding     string          `json:"funding" db:"funding"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type SearchFilters struct {
	Query   string          `form:"q"`
	Type    OpportunityType `form:"type"`
	Country string          `form:"country"`
	Page    int             `form:"page,default=1"`
	Limit   int             `form:"limit,default=20"`
}

type SearchResponse struct {
	Results    []Opportunity `json:"results"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
