package routes

import (
	"fmt"
	"log"
	"net/http"

	"silic0n-wiki/config"
	"silic0n-wiki/handlers"
	"silic0n-wiki/middleware"
)

func StartRouter() {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

	// Public routes
	mux.HandleFunc("GET /", handlers.Index)
	mux.HandleFunc("GET /wiki/{slug}", handlers.Article)
	mux.HandleFunc("GET /articles/recent", handlers.RecentArticles)
	mux.HandleFunc("GET /categories", handlers.Categories)
	mux.HandleFunc("GET /categories/{slug}", handlers.CategoryArticles)
	mux.HandleFunc("GET /categories/{category}/tags/{tag}", handlers.TagArticles)
	mux.HandleFunc("GET /api/search", handlers.Search)
	mux.HandleFunc("GET /media/{filename}", handlers.ServeMedia)

	// Auth routes
	mux.HandleFunc("GET /register", handlers.RegisterPage)
	mux.HandleFunc("POST /register", handlers.RegisterSubmit)
	mux.HandleFunc("GET /login", handlers.LoginPage)
	mux.HandleFunc("POST /login", handlers.LoginSubmit)
	mux.HandleFunc("POST /logout", middleware.RequireCSRF(handlers.Logout))

	// Protected routes (require auth + CSRF on POST)
	mux.HandleFunc("GET /wiki/new", middleware.RequireAuth(handlers.CreateArticlePage))
	mux.HandleFunc("POST /wiki/new", middleware.RequireAuth(middleware.RequireCSRF(handlers.CreateArticleSubmit)))
	mux.HandleFunc("GET /wiki/{slug}/edit", middleware.RequireAuth(handlers.EditArticlePage))
	mux.HandleFunc("POST /wiki/{slug}/edit", middleware.RequireAuth(middleware.RequireCSRF(handlers.EditArticleSubmit)))
	mux.HandleFunc("POST /api/media/upload", middleware.RequireAuth(middleware.RequireCSRF(handlers.MediaUpload)))

	// Wrap entire mux with session loading middleware
	wrappedMux := middleware.LoadSession(mux)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	err := http.ListenAndServe(addr, wrappedMux)
	if err != nil {
		log.Fatal(err)
	}
}
