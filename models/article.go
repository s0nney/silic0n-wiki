package models

import (
	"fmt"
	"strings"
	"time"

	"silic0n-wiki/database"
)

type Article struct {
	ID           int
	Slug         string
	Title        string
	Content      string
	CategoryID   int
	LastEditedBy string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ArticleWithCategory struct {
	Article
	CategoryName string
	CategorySlug string
	Tags         []TagWithCategory
}

func GetArticleBySlug(slug string) (*ArticleWithCategory, error) {
	article := &ArticleWithCategory{}
	err := database.DB.QueryRow(
		`SELECT a.id, a.slug, a.title, a.content, COALESCE(a.category_id, 0),
		        a.last_edited_by, a.created_at, a.updated_at,
		        COALESCE(c.name, '') as category_name, COALESCE(c.slug, '') as category_slug
		FROM articles a
		LEFT JOIN categories c ON a.category_id = c.id
		WHERE a.slug = $1`,
		slug,
	).Scan(&article.ID, &article.Slug, &article.Title, &article.Content, &article.CategoryID,
		&article.LastEditedBy, &article.CreatedAt, &article.UpdatedAt,
		&article.CategoryName, &article.CategorySlug)

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

func GetRecentArticles(limit int) ([]ArticleWithCategory, error) {
	rows, err := database.DB.Query(
		`SELECT a.id, a.slug, a.title, a.content, a.last_edited_by, a.created_at, a.updated_at,
		        COALESCE(c.name, '') as category_name, COALESCE(c.slug, '') as category_slug
		FROM articles a
		LEFT JOIN categories c ON a.category_id = c.id
		WHERE a.created_at >= NOW() - INTERVAL '48 hours'
		   OR a.updated_at >= NOW() - INTERVAL '48 hours'
		ORDER BY GREATEST(a.created_at, a.updated_at) DESC
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []ArticleWithCategory
	for rows.Next() {
		var a ArticleWithCategory
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.LastEditedBy,
			&a.CreatedAt, &a.UpdatedAt, &a.CategoryName, &a.CategorySlug); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, rows.Err()
}

func Slugify(title string) string {
	var result []byte
	prevHyphen := false
	for _, r := range strings.ToLower(strings.TrimSpace(title)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result = append(result, byte(r))
			prevHyphen = false
		} else if r == ' ' || r == '-' || r == '_' {
			if !prevHyphen && len(result) > 0 {
				result = append(result, '-')
				prevHyphen = true
			}
		}
	}
	if len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}
	return string(result)
}

func SlugExists(slug string, excludeID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM articles WHERE slug = $1 AND id != $2`,
		slug, excludeID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GenerateUniqueSlug(title string, excludeID int) (string, error) {
	base := Slugify(title)
	if base == "" {
		base = "article"
	}
	slug := base
	counter := 2
	for {
		exists, err := SlugExists(slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}
}

func CreateArticle(title, content string, categoryID int, lastEditedBy string) (*Article, error) {
	slug, err := GenerateUniqueSlug(title, 0)
	if err != nil {
		return nil, err
	}

	article := &Article{}
	err = database.DB.QueryRow(
		`INSERT INTO articles (slug, title, content, category_id, last_edited_by)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, slug, title, content, COALESCE(category_id, 0), last_edited_by, created_at, updated_at`,
		slug, title, content, categoryID, lastEditedBy,
	).Scan(&article.ID, &article.Slug, &article.Title, &article.Content,
		&article.CategoryID, &article.LastEditedBy, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func UpdateArticle(id int, title, content string, categoryID int, lastEditedBy string) (*Article, error) {
	slug, err := GenerateUniqueSlug(title, id)
	if err != nil {
		return nil, err
	}

	article := &Article{}
	err = database.DB.QueryRow(
		`UPDATE articles
		 SET slug = $1, title = $2, content = $3, category_id = $4,
		     last_edited_by = $5, updated_at = NOW()
		 WHERE id = $6
		 RETURNING id, slug, title, content, COALESCE(category_id, 0), last_edited_by, created_at, updated_at`,
		slug, title, content, categoryID, lastEditedBy, id,
	).Scan(&article.ID, &article.Slug, &article.Title, &article.Content,
		&article.CategoryID, &article.LastEditedBy, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func SetArticleTags(articleID int, tagIDs []int) error {
	_, err := database.DB.Exec(`DELETE FROM article_tags WHERE article_id = $1`, articleID)
	if err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		_, err := database.DB.Exec(
			`INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			articleID, tagID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
