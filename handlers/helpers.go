package handlers

import (
	"html/template"
	"log"
	"net/http"

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
	ts, err := template.ParseFiles(templates...)
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
