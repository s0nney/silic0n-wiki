package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"silic0n-wiki/middleware"
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

	tags, err := models.GetTagsForArticle(article.ID)
	if err != nil {
		log.Printf("Error fetching tags for article: %v", err)
	} else {
		article.Tags = tags
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/article.tmpl.html",
	}

	renderTemplate(w, r, files, article)
}

func CreateArticlePage(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/article_form.tmpl.html",
	}

	data := struct {
		IsEdit          bool
		Article         *models.Article
		Categories      []models.Category
		ArticleTags     []models.TagWithCategory
		TagString       string
		NewCategoryName string
		Errors          []string
	}{
		IsEdit:     false,
		Categories: categories,
	}

	renderTemplate(w, r, files, data)
}

func CreateArticleSubmit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	r.ParseForm()
	title := strings.TrimSpace(r.FormValue("title"))
	content := r.FormValue("content")
	categoryIDStr := r.FormValue("category_id")
	newCategoryName := strings.TrimSpace(r.FormValue("new_category_name"))
	tagsStr := r.FormValue("tags")

	var errors []string
	if title == "" {
		errors = append(errors, "Title is required")
	}
	if content == "" {
		errors = append(errors, "Content is required")
	}

	var categoryID int
	if categoryIDStr == "new" {
		if newCategoryName == "" {
			errors = append(errors, "Category name is required when creating a new category")
		}
	} else {
		categoryID, _ = strconv.Atoi(categoryIDStr)
		if categoryID == 0 {
			errors = append(errors, "Category is required")
		}
	}

	if len(errors) > 0 {
		categories, _ := models.GetAllCategories()
		files := []string{
			"./templates/base.tmpl.html",
			"./templates/article_form.tmpl.html",
		}
		data := struct {
			IsEdit          bool
			Article         *models.Article
			Categories      []models.Category
			ArticleTags     []models.TagWithCategory
			TagString       string
			NewCategoryName string
			Errors          []string
		}{
			IsEdit:          false,
			Article:         &models.Article{Title: title, Content: content, CategoryID: categoryID},
			Categories:      categories,
			TagString:       tagsStr,
			NewCategoryName: newCategoryName,
			Errors:          errors,
		}
		renderTemplate(w, r, files, data)
		return
	}

	if categoryIDStr == "new" {
		cat, err := models.GetOrCreateCategory(newCategoryName)
		if err != nil {
			log.Printf("Error creating category: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		categoryID = cat.ID
	}

	article, err := models.CreateArticle(title, content, categoryID, user.Username)
	if err != nil {
		log.Printf("Error creating article: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if tagsStr != "" {
		tagIDs, err := models.ResolveTagNames(tagsStr, categoryID)
		if err != nil {
			log.Printf("Error resolving tags: %v", err)
		} else if len(tagIDs) > 0 {
			if err := models.SetArticleTags(article.ID, tagIDs); err != nil {
				log.Printf("Error setting tags: %v", err)
			}
		}
	}

	http.Redirect(w, r, "/wiki/"+article.Slug, http.StatusSeeOther)
}

func EditArticlePage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
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

	categories, err := models.GetAllCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articleTags, err := models.GetTagsForArticle(article.ID)
	if err != nil {
		log.Printf("Error fetching article tags: %v", err)
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/article_form.tmpl.html",
	}

	formArticle := &models.Article{
		ID:           article.ID,
		Slug:         article.Slug,
		Title:        article.Title,
		Content:      article.Content,
		CategoryID:   article.CategoryID,
		LastEditedBy: article.LastEditedBy,
	}

	var tagNames []string
	for _, t := range articleTags {
		tagNames = append(tagNames, t.Name)
	}
	tagString := strings.Join(tagNames, ", ")

	data := struct {
		IsEdit          bool
		Article         *models.Article
		Categories      []models.Category
		ArticleTags     []models.TagWithCategory
		TagString       string
		NewCategoryName string
		Errors          []string
	}{
		IsEdit:      true,
		Article:     formArticle,
		Categories:  categories,
		ArticleTags: articleTags,
		TagString:   tagString,
	}

	renderTemplate(w, r, files, data)
}

func EditArticleSubmit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	slug := r.PathValue("slug")

	existingArticle, err := models.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching article: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	title := strings.TrimSpace(r.FormValue("title"))
	content := r.FormValue("content")
	categoryIDStr := r.FormValue("category_id")
	newCategoryName := strings.TrimSpace(r.FormValue("new_category_name"))
	tagsStr := r.FormValue("tags")

	var errors []string
	if title == "" {
		errors = append(errors, "Title is required")
	}
	if content == "" {
		errors = append(errors, "Content is required")
	}

	var categoryID int
	if categoryIDStr == "new" {
		if newCategoryName == "" {
			errors = append(errors, "Category name is required when creating a new category")
		}
	} else {
		categoryID, _ = strconv.Atoi(categoryIDStr)
		if categoryID == 0 {
			errors = append(errors, "Category is required")
		}
	}

	if len(errors) > 0 {
		categories, _ := models.GetAllCategories()
		articleTags, _ := models.GetTagsForArticle(existingArticle.ID)
		files := []string{
			"./templates/base.tmpl.html",
			"./templates/article_form.tmpl.html",
		}
		data := struct {
			IsEdit          bool
			Article         *models.Article
			Categories      []models.Category
			ArticleTags     []models.TagWithCategory
			TagString       string
			NewCategoryName string
			Errors          []string
		}{
			IsEdit:          true,
			Article:         &models.Article{ID: existingArticle.ID, Slug: slug, Title: title, Content: content, CategoryID: categoryID},
			Categories:      categories,
			ArticleTags:     articleTags,
			TagString:       tagsStr,
			NewCategoryName: newCategoryName,
			Errors:          errors,
		}
		renderTemplate(w, r, files, data)
		return
	}

	if categoryIDStr == "new" {
		cat, err := models.GetOrCreateCategory(newCategoryName)
		if err != nil {
			log.Printf("Error creating category: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		categoryID = cat.ID
	}

	updatedArticle, err := models.UpdateArticle(existingArticle.ID, title, content, categoryID, user.Username)
	if err != nil {
		log.Printf("Error updating article: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if tagsStr != "" {
		tagIDs, err := models.ResolveTagNames(tagsStr, categoryID)
		if err != nil {
			log.Printf("Error resolving tags: %v", err)
		} else {
			if err := models.SetArticleTags(updatedArticle.ID, tagIDs); err != nil {
				log.Printf("Error setting tags: %v", err)
			}
		}
	} else {
		if err := models.SetArticleTags(updatedArticle.ID, nil); err != nil {
			log.Printf("Error clearing tags: %v", err)
		}
	}

	http.Redirect(w, r, "/wiki/"+updatedArticle.Slug, http.StatusSeeOther)
}
