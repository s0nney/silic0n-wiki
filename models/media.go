package models

import (
	"time"

	"silic0n-wiki/database"
)

type Media struct {
	ID           int
	ArticleID    *int
	Filename     string
	OriginalName string
	FilePath     string
	MimeType     string
	FileSize     int64
	UploadedBy   string
	CreatedAt    time.Time
}

func CreateMedia(articleID *int, filename, originalName, filePath, mimeType string, fileSize int64, uploadedBy string) (*Media, error) {
	media := &Media{}
	err := database.DB.QueryRow(
		`INSERT INTO media (article_id, filename, original_name, file_path, mime_type, file_size, uploaded_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, article_id, filename, original_name, file_path, mime_type, file_size, uploaded_by, created_at`,
		articleID, filename, originalName, filePath, mimeType, fileSize, uploadedBy,
	).Scan(&media.ID, &media.ArticleID, &media.Filename, &media.OriginalName,
		&media.FilePath, &media.MimeType, &media.FileSize, &media.UploadedBy, &media.CreatedAt)
	if err != nil {
		return nil, err
	}
	return media, nil
}

func GetMediaByFilename(filename string) (*Media, error) {
	media := &Media{}
	err := database.DB.QueryRow(
		`SELECT id, article_id, filename, original_name, file_path, mime_type, file_size, uploaded_by, created_at
		 FROM media WHERE filename = $1`,
		filename,
	).Scan(&media.ID, &media.ArticleID, &media.Filename, &media.OriginalName,
		&media.FilePath, &media.MimeType, &media.FileSize, &media.UploadedBy, &media.CreatedAt)
	if err != nil {
		return nil, err
	}
	return media, nil
}

func GetMediaByID(id int) (*Media, error) {
	media := &Media{}
	err := database.DB.QueryRow(
		`SELECT id, article_id, filename, original_name, file_path, mime_type, file_size, uploaded_by, created_at
		 FROM media WHERE id = $1`,
		id,
	).Scan(&media.ID, &media.ArticleID, &media.Filename, &media.OriginalName,
		&media.FilePath, &media.MimeType, &media.FileSize, &media.UploadedBy, &media.CreatedAt)
	if err != nil {
		return nil, err
	}
	return media, nil
}

func GetMediaForArticle(articleID int) ([]Media, error) {
	rows, err := database.DB.Query(
		`SELECT id, article_id, filename, original_name, file_path, mime_type, file_size, uploaded_by, created_at
		 FROM media WHERE article_id = $1 ORDER BY created_at DESC`,
		articleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []Media
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.ArticleID, &m.Filename, &m.OriginalName,
			&m.FilePath, &m.MimeType, &m.FileSize, &m.UploadedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		mediaList = append(mediaList, m)
	}
	return mediaList, rows.Err()
}
