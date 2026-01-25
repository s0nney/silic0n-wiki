package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"silic0n-wiki/models"
)

type SearchResult struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	articles, err := models.SearchArticles(query)
	if err != nil {
		log.Printf("Error searching articles: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	results := make([]SearchResult, len(articles))
	for i, a := range articles {
		results[i] = SearchResult{
			Slug:        a.Slug,
			Title:       a.Title,
			Description: truncate(a.Content, 100),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
