package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

var ErrInvalidSession = errors.New("invalid session")

const sessionDuration = 12 * time.Hour

type SessionStore interface {
	Create(userID int, ip, userAgent string) (*models.Session, error)
	Validate(sessionID string) (*models.User, error)
	Revoke(sessionID string) error
}

type sessionStore struct {
	db *sql.DB
}

func NewSessionStore() SessionStore {
	return &sessionStore{db: db.DB}
}

func (s *sessionStore) Create(userID int, ip, userAgent string) (*models.Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(sessionDuration)

	_, err = s.db.Exec(`
		INSERT INTO sessions (id, user_id, expires_at, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?)`,
		sessionID, userID, expiry, ip, userAgent,
	)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiry,
	}, nil
}

func (s *sessionStore) Validate(sessionID string) (*models.User, error) {
	var u models.User
	err := s.db.QueryRow(`
		SELECT u.id, u.email, u.name, u.created_at, u.updated_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = ?
		  AND s.revoked_at IS NULL
		  AND s.expires_at > ?`,
		sessionID, time.Now(),
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *sessionStore) Revoke(sessionID string) error {
	result, err := s.db.Exec(`
		UPDATE sessions
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE id = ? AND revoked_at IS NULL`,
		sessionID,
	)
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
