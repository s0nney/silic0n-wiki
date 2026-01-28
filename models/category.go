package models

import (
	"fmt"
	"time"

	"silic0n-wiki/database"
)

type Category struct {
	ID          int
	Slug        string
	Name        string
	Description string
	CreatedAt   time.Time
}

type CategoryWithCount struct {
	Category
	ArticleCount int
}

func GetAllCategories() ([]Category, error) {
	rows, err := database.DB.Query(
		"SELECT id, slug, name, description, created_at FROM categories ORDER BY name",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func GetCategoriesWithArticleCount() ([]CategoryWithCount, error) {
	rows, err := database.DB.Query(`
		SELECT c.id, c.slug, c.name, c.description, c.created_at, COUNT(a.id) as article_count
		FROM categories c
		LEFT JOIN articles a ON a.category_id = c.id
		GROUP BY c.id, c.slug, c.name, c.description, c.created_at
		ORDER BY c.name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryWithCount
	for rows.Next() {
		var c CategoryWithCount
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.CreatedAt, &c.ArticleCount); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func GetCategoryBySlug(slug string) (*Category, error) {
	category := &Category{}
	err := database.DB.QueryRow(
		"SELECT id, slug, name, description, created_at FROM categories WHERE slug = $1",
		slug,
	).Scan(&category.ID, &category.Slug, &category.Name, &category.Description, &category.CreatedAt)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func GetOrCreateCategory(name string) (*Category, error) {
	slug := Slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("invalid category name")
	}

	category := &Category{}
	err := database.DB.QueryRow(
		`INSERT INTO categories (slug, name, description)
		 VALUES ($1, $2, '')
		 ON CONFLICT (slug) DO UPDATE SET slug = EXCLUDED.slug
		 RETURNING id, slug, name, description, created_at`,
		slug, name,
	).Scan(&category.ID, &category.Slug, &category.Name, &category.Description, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func GetArticlesByCategory(categoryID int) ([]Article, error) {
	rows, err := database.DB.Query(
		`SELECT id, slug, title, content, created_at, updated_at
		FROM articles
		WHERE category_id = $1
		ORDER BY title`,
		categoryID,
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
