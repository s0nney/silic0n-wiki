package models

import (
	"time"

	"silic0n-wiki/database"
)

type Session struct {
	Token     string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func CreateSession(token string, userID int, duration time.Duration) (*Session, error) {
	session := &Session{}
	expiresAt := time.Now().Add(duration)
	err := database.DB.QueryRow(
		`INSERT INTO sessions (token, user_id, expires_at)
		 VALUES ($1, $2, $3)
		 RETURNING token, user_id, created_at, expires_at`,
		token, userID, expiresAt,
	).Scan(&session.Token, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func GetSessionByToken(token string) (*Session, error) {
	session := &Session{}
	err := database.DB.QueryRow(
		`SELECT token, user_id, created_at, expires_at
		 FROM sessions
		 WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&session.Token, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func DeleteSession(token string) error {
	_, err := database.DB.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func DeleteExpiredSessions() error {
	_, err := database.DB.Exec(`DELETE FROM sessions WHERE expires_at < NOW()`)
	return err
}
