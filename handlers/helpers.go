package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"silic0n-wiki/auth"
	"silic0n-wiki/middleware"
	"silic0n-wiki/models"
)

type PageData struct {
	User      *models.User
	CSRFToken string
	Data      interface{}
}

func newPageData(r *http.Request, data interface{}) PageData {
	user := middleware.GetUser(r)
	sessionToken := middleware.GetSessionToken(r)
	csrfToken := ""
	if sessionToken != "" {
		csrfToken = auth.GenerateCSRFToken(sessionToken)
	}
	return PageData{
		User:      user,
		CSRFToken: csrfToken,
		Data:      data,
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, templates []string, data interface{}) {
	funcMap := template.FuncMap{
		"renderContent": RenderArticleContent,
	}

	ts, err := template.New("").Funcs(funcMap).ParseFiles(templates...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageData := newPageData(r, data)
	err = ts.ExecuteTemplate(w, "base.tmpl.html", pageData)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Matches ![alt](filename =WIDTHxHEIGHT center) — size and center are both optional
var mediaEmbedRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^\s)]+)(?:\s*=(\d*)x(\d*))?(\s+center)?\)`)

func RenderArticleContent(content string) template.HTML {
	escaped := template.HTMLEscapeString(content)

	result := mediaEmbedRegex.ReplaceAllStringFunc(escaped, func(match string) string {
		submatches := mediaEmbedRegex.FindStringSubmatch(match)
		if len(submatches) != 6 {
			return match
		}
		alt := submatches[1]
		filename := submatches[2]
		widthStr := submatches[3]
		heightStr := submatches[4]
		centered := strings.TrimSpace(submatches[5]) == "center"

		if !isValidMediaFilename(filename) {
			return match
		}

		mediaURL := "/media/" + filename
		ext := strings.ToLower(filepath.Ext(filename))

		style := buildSizeStyle(widthStr, heightStr)
		divClass := "media-embed"
		if centered {
			divClass += " media-center"
		}

		switch ext {
		case ".mp4", ".webm":
			return fmt.Sprintf(
				`<div class="%s media-video"><video controls preload="metadata" title="%s"%s><source src="%s">Your browser does not support video playback.</video></div>`,
				divClass, alt, style, mediaURL,
			)
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			return fmt.Sprintf(
				`<div class="%s media-image"><img src="%s" alt="%s"%s loading="lazy"></div>`,
				divClass, mediaURL, alt, style,
			)
		default:
			return match
		}
	})

	result = strings.ReplaceAll(result, "\n", "<br>\n")

	return template.HTML(result)
}

func buildSizeStyle(widthStr, heightStr string) string {
	if widthStr == "" && heightStr == "" {
		return ""
	}
	var parts []string
	if widthStr != "" {
		parts = append(parts, "width:"+widthStr+"px")
	}
	if heightStr != "" {
		parts = append(parts, "height:"+heightStr+"px")
	}
	return fmt.Sprintf(` style="%s"`, strings.Join(parts, ";"))
}

func isValidMediaFilename(filename string) bool {
	for _, c := range filename {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}
	return len(filename) > 0 && !strings.Contains(filename, "..")
}
