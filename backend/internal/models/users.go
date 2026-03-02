package models

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type userWithPassword struct {
	User
	PasswordHash string
}

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidName     = errors.New("name required")
	ErrInvalidPassword = errors.New("password must be at least 10 characters")
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

const bcryptCost = 12

func GetUserByID(id int) (*User, error) {
	var u User
	err := db.DB.QueryRow(`
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = ?`,
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func CreateUser(email, password, name string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// Input validation
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	if len(password) < 10 {
		return nil, ErrInvalidPassword
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := db.DB.Exec(`
		INSERT INTO users (email, password_hash, name)
		VALUES (?, ?, ?)`,
		email, string(hash), name,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrEmailExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetUserByID(int(id))
}

func AuthenticateUser(email, password string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	var u userWithPassword
	err := db.DB.QueryRow(`
		SELECT id, email, name, created_at, updated_at, password_hash
		FROM users
		WHERE email = ?`,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PasswordHash,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &u.User, nil
}

func CountUsers() (int, error) {
	var count int
	if err := db.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	if !errors.As(err, &sqliteErr) {
		return false
	}
	return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
