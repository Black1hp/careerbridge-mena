package main

import (
	"log"
	"net/http"
	"os"

	"github.com/black1hp/careerbridge-mena/backend/internal/database"
	"github.com/black1hp/careerbridge-mena/backend/internal/handlers"
	"github.com/black1hp/careerbridge-mena/backend/internal/search"
	"github.com/gorilla/mux"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer database.Close()

	if err := search.Connect(); err != nil {
		log.Fatal("Elasticsearch connection failed:", err)
	}

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/search", handlers.SearchOpportunities).Methods("GET")
	api.HandleFunc("/opportunities/{id}", handlers.GetOpportunity).Methods("GET")
	api.HandleFunc("/countries", handlers.GetCountries).Methods("GET")
	api.HandleFunc("/ingest", handlers.IngestOpportunities).Methods("POST")
	api.HandleFunc("/healthz", handlers.HealthCheck).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
