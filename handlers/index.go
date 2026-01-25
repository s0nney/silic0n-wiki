package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("silic0n-wiki-server", "Silic0n Wiki")

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/index.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.tmpl.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
