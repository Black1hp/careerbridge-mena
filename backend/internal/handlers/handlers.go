package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/black1hp/careerbridge-mena/backend/internal/database"
	"github.com/black1hp/careerbridge-mena/backend/internal/models"
	"github.com/black1hp/careerbridge-mena/backend/internal/search"
	"github.com/gorilla/mux"
)

func SearchOpportunities(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	req := search.SearchRequest{
		Query:   r.URL.Query().Get("q"),
		Type:    r.URL.Query().Get("type"),
		Country: r.URL.Query().Get("country"),
		Page:    page,
		Limit:   limit,
	}

	result, err := search.Search(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := result.Total / limit
	if result.Total%limit > 0 {
		totalPages++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SearchResponse{
		Results:    convertToModels(result.Opportunities),
		Total:      result.Total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

func GetOpportunity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var opp models.Opportunity
	err := database.Pool.QueryRow(r.Context(),
		`SELECT id, title, type, country, deadline, url, source, description, eligibility, funding, created_at, updated_at
		FROM opportunities WHERE id = $1`, id).
		Scan(&opp.ID, &opp.Title, &opp.Type, &opp.Country, &opp.Deadline, &opp.URL, &opp.Source, &opp.Description, &opp.Eligibility, &opp.Funding, &opp.CreatedAt, &opp.UpdatedAt)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opp)
}

func GetCountries(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(r.Context(), `SELECT DISTINCT country FROM opportunities ORDER BY country`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var countries []string
	for rows.Next() {
		var c string
		rows.Scan(&c)
		countries = append(countries, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

type IngestRequest struct {
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

func IngestOpportunities(w http.ResponseWriter, r *http.Request) {
	var items []IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var inserted int
	for _, item := range items {
		var deadline *time.Time
		if item.Deadline != nil {
			t, err := time.Parse("2006-01-02", *item.Deadline)
			if err == nil {
				deadline = &t
			}
		}

		var id int
		err := database.Pool.QueryRow(r.Context(),
			`INSERT INTO opportunities (title, type, country, deadline, url, source, description, eligibility, funding)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (title, source, url) DO UPDATE SET
				description = EXCLUDED.description,
				deadline = EXCLUDED.deadline,
				updated_at = NOW()
			RETURNING id`,
			item.Title, item.Type, item.Country, deadline, item.URL, item.Source, item.Description, item.Eligibility, item.Funding).
			Scan(&id)
		if err != nil {
			continue
		}

		search.IndexOpportunity(r.Context(), search.ESOpportunity{
			ID:          id,
			Title:       item.Title,
			Type:        item.Type,
			Country:     item.Country,
			URL:         item.URL,
			Source:      item.Source,
			Description: item.Description,
			Eligibility: item.Eligibility,
			Funding:     item.Funding,
		})
		inserted++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"inserted": inserted})
}

func convertToModels(opps []search.ESOpportunity) []models.Opportunity {
	var result []models.Opportunity
	for _, o := range opps {
		var opp models.Opportunity
		opp.ID = o.ID
		opp.Title = o.Title
		opp.Type = models.OpportunityType(o.Type)
		opp.Country = o.Country
		opp.URL = o.URL
		opp.Source = o.Source
		opp.Description = o.Description
		opp.Eligibility = o.Eligibility
		opp.Funding = o.Funding
		if o.Deadline != "" {
			t, err := time.Parse("2006-01-02", o.Deadline[:10])
			if err == nil {
				opp.Deadline = &t
			}
		}
		result = append(result, opp)
	}
	return result
}
