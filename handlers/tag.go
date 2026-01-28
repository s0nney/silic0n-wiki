package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"silic0n-wiki/models"
)

func TagArticles(w http.ResponseWriter, r *http.Request) {
	categorySlug := r.PathValue("category")
	tagSlug := r.PathValue("tag")

	if categorySlug == "" || tagSlug == "" {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	tag, err := models.GetTagBySlug(tagSlug, categorySlug)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tag not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching tag: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articles, err := models.GetArticlesByTag(tag.ID)
	if err != nil {
		log.Printf("Error fetching articles for tag: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/tag.tmpl.html",
	}

	data := struct {
		Tag      *models.TagWithCategory
		Articles []models.Article
	}{
		Tag:      tag,
		Articles: articles,
	}

	renderTemplate(w, r, files, data)
}
