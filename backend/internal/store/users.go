package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const bcryptCost = 12

type UserStore interface {
	Count() (int, error)
	GetByID(id int) (*models.User, error)
	Create(email, password, name string) (*models.User, error)
	Authenticate(email, password string) (*models.User, error)
}

type userStore struct {
	db *sql.DB
}

func NewUserStore() UserStore {
	return &userStore{db: db.DB}
}

type userWithPassword struct {
	models.User
	PasswordHash string
}

func (s *userStore) Count() (int, error) {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *userStore) GetByID(id int) (*models.User, error) {
	var u models.User
	err := s.db.QueryRow(`
		SELECT id, email, name, created_at, updated_at
		FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *userStore) Create(email, password, name string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := s.db.Exec(`
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

	return s.GetByID(int(id))
}

func (s *userStore) Authenticate(email, password string) (*models.User, error) {
	var u userWithPassword
	err := s.db.QueryRow(`
		SELECT id, email, name, created_at, updated_at, password_hash
		FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt, &u.PasswordHash)
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
