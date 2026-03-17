package services

import (
	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

type StatusService interface {
	List(projectID int) ([]models.Status, error)
	Create(projectID int, label, color string) (*models.Status, error)
	Update(statusID int, label, color *string) (*models.Status, error)
	UpdatePosition(statusID, position int) (*models.Status, error)
	Delete(statusID int) error
}

type statusService struct {
	store    store.StatusStore
	projects store.ProjectStore
}

func NewStatusService(store store.StatusStore, projects store.ProjectStore) StatusService {
	return &statusService{store: store, projects: projects}
}

func (s *statusService) List(projectID int) ([]models.Status, error) {
	return s.store.List(projectID)
}

func (s *statusService) Create(projectID int, label, color string) (*models.Status, error) {
	if label == "" {
		return nil, ErrLabelRequired
	}
	if color == "" {
		color = defaultColor
	} else if !colorRegex.MatchString(color) {
		return nil, ErrInvalidColor
	}
	if _, err := s.projects.GetByID(projectID); err != nil {
		return nil, err
	}
	return s.store.Create(projectID, label, color)
}

func (s *statusService) Update(statusID int, label, color *string) (*models.Status, error) {
	if color != nil && !colorRegex.MatchString(*color) {
		return nil, ErrInvalidColor
	}
	return s.store.Update(statusID, label, color)
}

func (s *statusService) UpdatePosition(statusID, position int) (*models.Status, error) {
	return s.store.UpdatePosition(statusID, position)
}

func (s *statusService) Delete(statusID int) error {
	return s.store.Delete(statusID)
}
