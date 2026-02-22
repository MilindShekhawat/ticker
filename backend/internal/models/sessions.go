package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID int, ip, userAgent string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(12 * time.Hour)

	query := `
		INSERT INTO sessions (id, user_id, expires_at, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.DB.Exec(query, sessionID, userID, expiry, ip, userAgent)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiry,
	}, nil
}

func ValidateSession(sessionID string) (*User, error) {
	query := `
		SELECT users.id, users.email, users.name,
		       users.created_at, users.updated_at
		FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.id = ?
		  AND sessions.revoked_at IS NULL
		  AND sessions.expires_at > CURRENT_TIMESTAMP
	`

	var u User
	err := db.DB.QueryRow(query, sessionID).Scan(
		&u.ID, &u.Email, &u.Name,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	return &u, nil
}

func RevokeSession(sessionID string) error {
	query := `
		UPDATE sessions
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.DB.Exec(query, sessionID)
	return err
}
