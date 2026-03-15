package services

import (
	"strings"
	"unicode"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

type ProjectService interface {
	ListForUser(userID int) ([]models.Project, error)
	GetForUser(projectID, userID int) (*models.Project, error)
	CreateForUser(name, description, keyPrefix string, userID int) (*models.Project, error)
	UpdateForUser(id, userID int, name, description *string) (*models.Project, error)
	DeleteForUser(id, userID int) error
}

type projectService struct {
	store store.ProjectStore
}

func NewProjectService(store store.ProjectStore) ProjectService {
	return &projectService{store: store}
}

func (s *projectService) ListForUser(userID int) ([]models.Project, error) {
	return s.store.ListByUser(userID)
}

func (s *projectService) GetForUser(projectID, userID int) (*models.Project, error) {
	return s.store.GetByUser(projectID, userID)
}

func (s *projectService) CreateForUser(name, description, keyPrefix string, userID int) (*models.Project, error) {
	name = strings.TrimSpace(name)
	keyPrefix = strings.ToUpper(strings.TrimSpace(keyPrefix))

	if name == "" || len(name) > models.MaxProjectNameLength {
		return nil, ErrInvalidProjectName
	}
	if !isValidKeyPrefix(keyPrefix) {
		return nil, ErrInvalidKeyPrefix
	}
	if userID == 0 {
		return nil, ErrInvalidProject
	}

	return s.store.Create(name, description, keyPrefix, userID)
}

func (s *projectService) UpdateForUser(id, userID int, name, description *string) (*models.Project, error) {
	if name == nil && description == nil {
		return nil, ErrInvalidProjectUpdate
	}
	if name != nil {
		n := strings.TrimSpace(*name)
		if n == "" || len(n) > models.MaxProjectNameLength {
			return nil, ErrInvalidProjectName
		}
		name = &n
	}

	return s.store.UpdateByUser(id, userID, name, description)
}

func (s *projectService) DeleteForUser(id, userID int) error {
	return s.store.DeleteByUser(id, userID)
}

func isValidKeyPrefix(s string) bool {
	if s == "" || len(s) > 5 {
		return false
	}
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
