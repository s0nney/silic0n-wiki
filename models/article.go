package models

import (
	"time"

	"silic0n-wiki/database"
)

type Article struct {
	ID        int
	Slug      string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetArticleBySlug(slug string) (*Article, error) {
	article := &Article{}
	err := database.DB.QueryRow(
		"SELECT id, slug, title, content, created_at, updated_at FROM articles WHERE slug = $1",
		slug,
	).Scan(&article.ID, &article.Slug, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return article, nil
}

func SearchArticles(query string) ([]Article, error) {
	rows, err := database.DB.Query(
		`SELECT id, slug, title, content, created_at, updated_at
		FROM articles
		WHERE title ILIKE '%' || $1 || '%' OR slug ILIKE '%' || $1 || '%'
		ORDER BY title
		LIMIT 10`,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, rows.Err()
}

func GetAllArticles() ([]Article, error) {
	rows, err := database.DB.Query(
		"SELECT id, slug, title, content, created_at, updated_at FROM articles ORDER BY title",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, rows.Err()
}

func GetRecentArticles(limit int) ([]Article, error) {
	rows, err := database.DB.Query(
		`SELECT id, slug, title, content, created_at, updated_at
		FROM articles
		ORDER BY created_at DESC
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, rows.Err()
}
