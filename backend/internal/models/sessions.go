package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

var (
	ErrInvalidSession = errors.New("invalid session")
)

const sessionDuration = 12 * time.Hour

func CreateSession(userID int, ip, userAgent string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(sessionDuration)

	query := `
		INSERT INTO sessions (id, user_id, expires_at, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?)
	`

	if _, err := db.DB.Exec(query, sessionID, userID, expiry, ip, userAgent); err != nil {
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
		SELECT users.id, users.email, users.name, users.created_at, users.updated_at
		FROM sessions
		JOIN users ON users.id = sessions.user_id
		WHERE sessions.id = ?
		  AND sessions.revoked_at IS NULL
		  AND sessions.expires_at > ?
	`

	var u User
	err := db.DB.QueryRow(query, sessionID, time.Now()).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func RevokeSession(sessionID string) error {
	query := `
		UPDATE sessions
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND revoked_at IS NULL
	`

	result, err := db.DB.Exec(query, sessionID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrInvalidSession
	}

	return nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)                   // 256 bits
	if _, err := rand.Read(b); err != nil { // used crypto/rand
		return "", err
	}
	// convert byte to url safe string
	return base64.RawURLEncoding.EncodeToString(b), nil
}
