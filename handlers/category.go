package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"silic0n-wiki/models"
)

func Categories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetCategoriesWithArticleCount()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/categories.tmpl.html",
	}

	data := struct {
		Categories []models.CategoryWithCount
	}{
		Categories: categories,
	}

	renderTemplate(w, r, files, data)
}

func CategoryArticles(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category, err := models.GetCategoryBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching category: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articles, err := models.GetArticlesByCategory(category.ID)
	if err != nil {
		log.Printf("Error fetching articles for category: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tags, err := models.GetTagsByCategoryWithCount(category.ID)
	if err != nil {
		log.Printf("Error fetching tags for category: %v", err)
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/category.tmpl.html",
	}

	data := struct {
		Category *models.Category
		Articles []models.Article
		Tags     []models.TagWithCount
	}{
		Category: category,
		Articles: articles,
		Tags:     tags,
	}

	renderTemplate(w, r, files, data)
}
