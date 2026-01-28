package handlers

import (
	"log"
	"net/http"

	"silic0n-wiki/models"
)

func RecentArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := models.GetRecentArticles(10)
	if err != nil {
		log.Printf("Error fetching recent articles: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/recent.tmpl.html",
	}

	data := struct {
		Articles []models.ArticleWithCategory
	}{
		Articles: articles,
	}

	renderTemplate(w, r, files, data)
}
