package services

import (
	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

type PriorityService interface {
	List(projectID int) ([]models.Priority, error)
	Create(projectID int, label, color string) (*models.Priority, error)
	Update(priorityID int, label, color *string) (*models.Priority, error)
	UpdatePosition(priorityID, position int) (*models.Priority, error)
	Delete(priorityID int) error
}

type priorityService struct {
	store    store.PriorityStore
	projects store.ProjectStore
}

func NewPriorityService(store store.PriorityStore, projects store.ProjectStore) PriorityService {
	return &priorityService{store: store, projects: projects}
}

func (s *priorityService) List(projectID int) ([]models.Priority, error) {
	return s.store.List(projectID)
}

func (s *priorityService) Create(projectID int, label, color string) (*models.Priority, error) {
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

func (s *priorityService) Update(priorityID int, label, color *string) (*models.Priority, error) {
	if color != nil && !colorRegex.MatchString(*color) {
		return nil, ErrInvalidColor
	}
	return s.store.Update(priorityID, label, color)
}

func (s *priorityService) UpdatePosition(priorityID, position int) (*models.Priority, error) {
	return s.store.UpdatePosition(priorityID, position)
}

func (s *priorityService) Delete(priorityID int) error {
	return s.store.Delete(priorityID)
}
