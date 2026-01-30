package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"silic0n-wiki/config"
	"silic0n-wiki/middleware"
	"silic0n-wiki/models"
)

func MediaUpload(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	maxSize := config.AppConfig.Media.MaxFileSize
	r.Body = http.MaxBytesReader(w, r.Body, maxSize+1024)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		jsonError(w, "File too large. Maximum size is 10MB.", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		jsonError(w, "No file provided.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !isAllowedType(mimeType) {
		jsonError(w, "File type not allowed. Allowed: JPEG, PNG, GIF, WebP, MP4, WebM.", http.StatusBadRequest)
		return
	}

	if header.Size > maxSize {
		jsonError(w, "File too large. Maximum size is 10MB.", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = extensionFromMIME(mimeType)
	}
	uuidName := generateUUID() + ext
	filePath := filepath.Join(config.AppConfig.Media.UploadDir, uuidName)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		jsonError(w, "Failed to save file.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error writing file: %v", err)
		os.Remove(filePath)
		jsonError(w, "Failed to save file.", http.StatusInternalServerError)
		return
	}

	media, err := models.CreateMedia(nil, uuidName, header.Filename, filePath, mimeType, header.Size, user.Username)
	if err != nil {
		log.Printf("Error saving media record: %v", err)
		os.Remove(filePath)
		jsonError(w, "Failed to save media record.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            media.ID,
		"filename":      media.Filename,
		"original_name": media.OriginalName,
		"mime_type":     media.MimeType,
		"file_size":     media.FileSize,
		"embed_tag":     fmt.Sprintf("![%s](%s)", media.OriginalName, media.Filename),
		"preview_url":   fmt.Sprintf("/media/%s", media.Filename),
	})
}

func ServeMedia(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	if filename == "" {
		http.NotFound(w, r)
		return
	}

	filename = filepath.Base(filename)

	media, err := models.GetMediaByFilename(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(config.AppConfig.Media.UploadDir, media.Filename)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening media file: %v", err)
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", media.MimeType)
	w.Header().Set("Cache-Control", "public, max-age=86400")

	http.ServeContent(w, r, media.Filename, media.CreatedAt, file)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func isAllowedType(mimeType string) bool {
	for _, allowed := range config.AppConfig.Media.AllowedTypes {
		if strings.EqualFold(mimeType, allowed) {
			return true
		}
	}
	return false
}

func extensionFromMIME(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	default:
		return ""
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
