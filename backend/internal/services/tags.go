package services

import (
	"errors"
	"regexp"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

var (
	ErrLabelRequired = errors.New("label is required")
	ErrInvalidColor  = errors.New("invalid color format: must be #RRGGBB")
)

var colorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

const defaultColor = "#808080"

type TagService interface {
	List(projectID int) ([]models.Tag, error)
	Create(projectID int, label, color string) (*models.Tag, error)
	Update(tagID int, label, color *string) (*models.Tag, error)
	Delete(tagID int) error
}

type tagService struct {
	store    store.TagStore
	projects store.ProjectStore
}

func NewTagService(store store.TagStore, projects store.ProjectStore) TagService {
	return &tagService{store: store, projects: projects}
}

func (s *tagService) List(projectID int) ([]models.Tag, error) {
	return s.store.List(projectID)
}

func (s *tagService) Create(projectID int, label, color string) (*models.Tag, error) {
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

func (s *tagService) Update(tagID int, label, color *string) (*models.Tag, error) {
	if color != nil && !colorRegex.MatchString(*color) {
		return nil, ErrInvalidColor
	}

	return s.store.Update(tagID, label, color)
}

func (s *tagService) Delete(tagID int) error {
	return s.store.Delete(tagID)
}
