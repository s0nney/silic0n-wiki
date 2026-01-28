package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"silic0n-wiki/auth"
	"silic0n-wiki/middleware"
	"silic0n-wiki/models"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if middleware.GetUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/register.tmpl.html",
	}

	data := struct {
		Errors   []string
		Username string
		Email    string
	}{}

	renderTemplate(w, r, files, data)
}

func RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	if middleware.GetUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	var errors []string

	if len(username) < 3 || len(username) > 50 {
		errors = append(errors, "Username must be between 3 and 50 characters")
	} else if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username); !matched {
		errors = append(errors, "Username may only contain letters, numbers, and underscores")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		errors = append(errors, "Please enter a valid email address")
	}

	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters")
	}

	if password != passwordConfirm {
		errors = append(errors, "Passwords do not match")
	}

	if len(errors) == 0 {
		if _, err := models.GetUserByUsername(username); err == nil {
			errors = append(errors, "Username is already taken")
		} else if err != sql.ErrNoRows {
			log.Printf("Error checking username: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if _, err := models.GetUserByEmail(email); err == nil {
			errors = append(errors, "Email is already registered")
		} else if err != sql.ErrNoRows {
			log.Printf("Error checking email: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if len(errors) > 0 {
		files := []string{
			"./templates/base.tmpl.html",
			"./templates/register.tmpl.html",
		}
		data := struct {
			Errors   []string
			Username string
			Email    string
		}{
			Errors:   errors,
			Username: username,
			Email:    email,
		}
		renderTemplate(w, r, files, data)
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := models.CreateUser(username, email, hash)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := createSessionAndRedirect(w, r, user.ID, "/"); err != nil {
		log.Printf("Error creating session after registration: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if middleware.GetUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	files := []string{
		"./templates/base.tmpl.html",
		"./templates/login.tmpl.html",
	}

	data := struct {
		Error    string
		Username string
	}{}

	renderTemplate(w, r, files, data)
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if middleware.GetUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	user, err := models.GetUserByUsername(username)
	if err != nil {
		files := []string{
			"./templates/base.tmpl.html",
			"./templates/login.tmpl.html",
		}
		data := struct {
			Error    string
			Username string
		}{
			Error:    "Invalid username or password",
			Username: username,
		}
		renderTemplate(w, r, files, data)
		return
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		files := []string{
			"./templates/base.tmpl.html",
			"./templates/login.tmpl.html",
		}
		data := struct {
			Error    string
			Username string
		}{
			Error:    "Invalid username or password",
			Username: username,
		}
		renderTemplate(w, r, files, data)
		return
	}

	if err := createSessionAndRedirect(w, r, user.ID, "/"); err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionToken := middleware.GetSessionToken(r)
	if sessionToken != "" {
		models.DeleteSession(sessionToken)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createSessionAndRedirect(w http.ResponseWriter, r *http.Request, userID int, redirectTo string) error {
	rawToken, err := auth.GenerateToken(32)
	if err != nil {
		return err
	}

	_, err = models.CreateSession(rawToken, userID, 7*24*time.Hour)
	if err != nil {
		return err
	}

	signedToken := auth.SignToken(rawToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    signedToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	return nil
}
