package middleware

import (
	"context"
	"net/http"

	"silic0n-wiki/auth"
	"silic0n-wiki/models"
)

type contextKey string

const (
	UserContextKey    contextKey = "user"
	SessionContextKey contextKey = "session_token"
)

func GetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func GetSessionToken(r *http.Request) string {
	token, ok := r.Context().Value(SessionContextKey).(string)
	if !ok {
		return ""
	}
	return token
}

func LoadSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		rawToken, valid := auth.VerifySignedToken(cookie.Value)
		if !valid {
			next.ServeHTTP(w, r)
			return
		}

		session, err := models.GetSessionByToken(rawToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := models.GetUserByID(session.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, SessionContextKey, rawToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if GetUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken := GetSessionToken(r)
		if sessionToken == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		csrfToken := r.FormValue("csrf_token")
		if csrfToken == "" {
			csrfToken = r.Header.Get("X-CSRF-Token")
		}
		if !auth.ValidateCSRFToken(csrfToken, sessionToken) {
			http.Error(w, "Forbidden - invalid CSRF token", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
