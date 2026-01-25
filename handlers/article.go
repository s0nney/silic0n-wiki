package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"silic0n-wiki/models"
)

func Article(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	article, err := models.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching article: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/article.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.tmpl.html", article)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
