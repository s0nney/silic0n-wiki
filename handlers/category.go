package handlers

import (
	"database/sql"
	"html/template"
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

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []models.CategoryWithCount
	}{
		Categories: categories,
	}

	err = ts.ExecuteTemplate(w, "base.tmpl.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/category.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Category *models.Category
		Articles []models.Article
	}{
		Category: category,
		Articles: articles,
	}

	err = ts.ExecuteTemplate(w, "base.tmpl.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
