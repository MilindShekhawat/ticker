package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("password must be at least 10 characters")
	ErrInvalidName     = errors.New("name required")
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

const MinPasswordLength = 8

type AuthService interface {
	Signup(email, password, name, ip, userAgent string) (*models.User, *models.Session, error)
	Login(email, password, ip, userAgent string) (*models.User, *models.Session, error)
	Logout(sessionID string) error
}

type authService struct {
	users    store.UserStore
	sessions store.SessionStore
}

func NewAuthService(users store.UserStore, sessions store.SessionStore) AuthService {
	return &authService{users: users, sessions: sessions}
}

func (s *authService) Signup(email, password, name, ip, userAgent string) (*models.User, *models.Session, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !emailRegex.MatchString(email) {
		return nil, nil, ErrInvalidEmail
	}
	if len(password) < MinPasswordLength {
		return nil, nil, ErrInvalidPassword
	}
	if strings.TrimSpace(name) == "" {
		return nil, nil, ErrInvalidName
	}

	count, err := s.users.Count()
	if err != nil {
		return nil, nil, err
	}

	user, err := s.users.Create(email, password, name)
	if err != nil {
		return nil, nil, err
	}

	// first user ever → could assign admin role here
	_ = count

	session, err := s.sessions.Create(user.ID, ip, userAgent)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *authService) Login(email, password, ip, userAgent string) (*models.User, *models.Session, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.users.Authenticate(email, password)
	if err != nil {
		return nil, nil, err
	}

	session, err := s.sessions.Create(user.ID, ip, userAgent)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *authService) Logout(sessionID string) error {
	return s.sessions.Revoke(sessionID)
}
