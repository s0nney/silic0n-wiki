package routes

import (
	"fmt"
	"log"
	"net/http"

	"silic0n-wiki/config"
	"silic0n-wiki/handlers"
)

func StartRouter() {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /", handlers.Index)
	mux.HandleFunc("GET /wiki/{slug}", handlers.Article)
	mux.HandleFunc("GET /articles/recent", handlers.RecentArticles)
	mux.HandleFunc("GET /categories", handlers.Categories)
	mux.HandleFunc("GET /categories/{slug}", handlers.CategoryArticles)
	mux.HandleFunc("GET /api/search", handlers.Search)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
