package models

import (
	"fmt"
	"strings"
	"time"

	"silic0n-wiki/database"
)

type Tag struct {
	ID         int
	Slug       string
	Name       string
	CategoryID int
	CreatedAt  time.Time
}

type TagWithCategory struct {
	Tag
	CategoryName string
	CategorySlug string
}

type TagWithCount struct {
	Tag
	ArticleCount int
}

func GetTagBySlug(slug string, categorySlug string) (*TagWithCategory, error) {
	tag := &TagWithCategory{}
	err := database.DB.QueryRow(
		`SELECT t.id, t.slug, t.name, t.category_id, t.created_at,
		        c.name as category_name, c.slug as category_slug
		FROM tags t
		JOIN categories c ON t.category_id = c.id
		WHERE t.slug = $1 AND c.slug = $2`,
		slug, categorySlug,
	).Scan(&tag.ID, &tag.Slug, &tag.Name, &tag.CategoryID, &tag.CreatedAt,
		&tag.CategoryName, &tag.CategorySlug)

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func GetTagsByCategory(categoryID int) ([]Tag, error) {
	rows, err := database.DB.Query(
		`SELECT id, slug, name, category_id, created_at
		FROM tags
		WHERE category_id = $1
		ORDER BY name`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.CategoryID, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}

func GetTagsByCategoryWithCount(categoryID int) ([]TagWithCount, error) {
	rows, err := database.DB.Query(
		`SELECT t.id, t.slug, t.name, t.category_id, t.created_at, COUNT(at.article_id) as article_count
		FROM tags t
		LEFT JOIN article_tags at ON t.id = at.tag_id
		WHERE t.category_id = $1
		GROUP BY t.id, t.slug, t.name, t.category_id, t.created_at
		ORDER BY t.name`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []TagWithCount
	for rows.Next() {
		var t TagWithCount
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.CategoryID, &t.CreatedAt, &t.ArticleCount); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}

func GetTagsForArticle(articleID int) ([]TagWithCategory, error) {
	rows, err := database.DB.Query(
		`SELECT t.id, t.slug, t.name, t.category_id, t.created_at,
		        c.name as category_name, c.slug as category_slug
		FROM tags t
		JOIN article_tags at ON t.id = at.tag_id
		JOIN categories c ON t.category_id = c.id
		WHERE at.article_id = $1
		ORDER BY t.name`,
		articleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []TagWithCategory
	for rows.Next() {
		var t TagWithCategory
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.CategoryID, &t.CreatedAt,
			&t.CategoryName, &t.CategorySlug); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}

func GetOrCreateTag(name string, categoryID int) (*Tag, error) {
	slug := Slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("invalid tag name")
	}

	tag := &Tag{}
	err := database.DB.QueryRow(
		`INSERT INTO tags (slug, name, category_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (slug, category_id) DO UPDATE SET slug = EXCLUDED.slug
		 RETURNING id, slug, name, category_id, created_at`,
		slug, name, categoryID,
	).Scan(&tag.ID, &tag.Slug, &tag.Name, &tag.CategoryID, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func ResolveTagNames(commaSeparated string, categoryID int) ([]int, error) {
	var tagIDs []int
	for _, raw := range strings.Split(commaSeparated, ",") {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		tag, err := GetOrCreateTag(name, categoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve tag %q: %w", name, err)
		}
		tagIDs = append(tagIDs, tag.ID)
	}
	return tagIDs, nil
}

func GetArticlesByTag(tagID int) ([]Article, error) {
	rows, err := database.DB.Query(
		`SELECT a.id, a.slug, a.title, a.content, a.last_edited_by, a.created_at, a.updated_at
		FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		WHERE at.tag_id = $1
		ORDER BY a.title`,
		tagID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.LastEditedBy, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, rows.Err()
}
