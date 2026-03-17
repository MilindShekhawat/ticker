package services

import (
	"errors"

	"github.com/MilindShekhawat/ticker/internal/models"
	"github.com/MilindShekhawat/ticker/internal/store"
)

var (
	ErrInvalidRole             = errors.New("invalid role")
	ErrInsufficientProjectRole = errors.New("insufficient project role")
)

type ProjectMemberService interface {
	ListForUser(projectID, actorUserID int) ([]models.ProjectMember, error)
	AddForUser(projectID, actorUserID, userID, role int) error
	UpdateRoleForUser(projectID, actorUserID, userID, role int) error
	RemoveForUser(projectID, actorUserID, userID int) error
	TransferOwnershipForUser(projectID, actorUserID, newOwnerID int) error
}

type projectMemberService struct {
	store store.ProjectMemberStore
}

func NewProjectMemberService(store store.ProjectMemberStore) ProjectMemberService {
	return &projectMemberService{store: store}
}

func (s *projectMemberService) ListForUser(projectID, actorUserID int) ([]models.ProjectMember, error) {
	if _, err := s.store.GetUserRole(projectID, actorUserID); err != nil {
		return nil, err
	}

	return s.store.ListByProject(projectID)
}

func (s *projectMemberService) AddForUser(projectID, actorUserID, userID, role int) error {
	if !isValidRole(role) {
		return ErrInvalidRole
	}
	if role == models.RoleOwner {
		return store.ErrOwnerAlreadyExists
	}

	actorRole, err := s.store.GetUserRole(projectID, actorUserID)
	if err != nil {
		return err
	}
	if !isAdminOrOwner(actorRole) {
		return ErrInsufficientProjectRole
	}

	return s.store.Add(projectID, userID, role)
}

func (s *projectMemberService) UpdateRoleForUser(projectID, actorUserID, userID, role int) error {
	if !isValidRole(role) {
		return ErrInvalidRole
	}
	if role == models.RoleOwner {
		return store.ErrCannotPromoteToOwner
	}

	actorRole, err := s.store.GetUserRole(projectID, actorUserID)
	if err != nil {
		return err
	}
	if actorRole != models.RoleOwner {
		return ErrInsufficientProjectRole
	}

	targetRole, err := s.store.GetUserRole(projectID, userID)
	if err != nil {
		return err
	}
	if targetRole == models.RoleOwner {
		return store.ErrCannotDemoteOwner
	}

	return s.store.UpdateRole(projectID, userID, role)
}

func (s *projectMemberService) RemoveForUser(projectID, actorUserID, userID int) error {
	actorRole, err := s.store.GetUserRole(projectID, actorUserID)
	if err != nil {
		return err
	}
	if !isAdminOrOwner(actorRole) {
		return ErrInsufficientProjectRole
	}

	targetRole, err := s.store.GetUserRole(projectID, userID)
	if err != nil {
		return err
	}
	if targetRole == models.RoleOwner {
		return store.ErrCannotRemoveOwner
	}

	return s.store.Remove(projectID, userID)
}

func (s *projectMemberService) TransferOwnershipForUser(projectID, actorUserID, newOwnerID int) error {
	actorRole, err := s.store.GetUserRole(projectID, actorUserID)
	if err != nil {
		return err
	}
	if actorRole != models.RoleOwner {
		return ErrInsufficientProjectRole
	}

	return s.store.TransferOwnership(projectID, actorUserID, newOwnerID)
}

func isValidRole(role int) bool {
	return role >= models.RoleOwner && role <= models.RoleViewer
}

func isAdminOrOwner(role int) bool {
	return role == models.RoleOwner || role == models.RoleAdmin
}
