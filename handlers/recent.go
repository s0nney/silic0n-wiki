package handlers

import (
	"html/template"
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

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Articles []models.Article
	}{
		Articles: articles,
	}

	err = ts.ExecuteTemplate(w, "base.tmpl.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
