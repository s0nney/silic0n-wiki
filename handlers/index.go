package handlers

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("silic0n-wiki-server", "Silic0n Wiki")

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/index.tmpl.html",
	}

	renderTemplate(w, r, files, nil)
}
