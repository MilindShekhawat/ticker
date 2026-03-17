package store

import (
	"database/sql"
	"errors"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

var (
	ErrNotProjectMember      = errors.New("not a project member")
	ErrNotProjectOwner       = errors.New("not project owner")
	ErrCannotRemoveOwner     = errors.New("cannot remove project owner")
	ErrCannotDemoteOwner     = errors.New("cannot demote project owner")
	ErrMemberAlreadyExists   = errors.New("member already exists")
	ErrOwnerAlreadyExists    = errors.New("owner already exists")
	ErrCannotPromoteToOwner  = errors.New("use transfer ownership instead")
	ErrOwnershipTransferFail = errors.New("ownership transfer failed")
)

type ProjectMemberStore interface {
	GetUserRole(projectID, userID int) (int, error)
	ListByProject(projectID int) ([]models.ProjectMember, error)
	Add(projectID, userID, role int) error
	UpdateRole(projectID, userID, role int) error
	Remove(projectID, userID int) error
	TransferOwnership(projectID, currentOwnerID, newOwnerID int) error
}

type projectMemberStore struct {
	db *sql.DB
}

func NewProjectMemberStore() ProjectMemberStore {
	return &projectMemberStore{db: db.DB}
}

func (s *projectMemberStore) GetUserRole(projectID, userID int) (int, error) {
	var role int
	err := s.db.QueryRow(`
		SELECT role
		FROM project_members
		WHERE project_id = ? AND user_id = ?`,
		projectID,
		userID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotProjectMember
	}
	if err != nil {
		return 0, err
	}

	return role, nil
}

func (s *projectMemberStore) ListByProject(projectID int) ([]models.ProjectMember, error) {
	rows, err := s.db.Query(`
		SELECT project_id, user_id, role
		FROM project_members
		WHERE project_id = ?`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.ProjectMember, 0)
	for rows.Next() {
		var m models.ProjectMember
		if err := rows.Scan(&m.ProjectID, &m.UserID, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, rows.Err()
}

func (s *projectMemberStore) Add(projectID, userID, role int) error {
	_, err := s.db.Exec(`
		INSERT INTO project_members (project_id, user_id, role)
		VALUES (?, ?, ?)`,
		projectID,
		userID,
		role,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrMemberAlreadyExists
		}
		return err
	}

	return nil
}

func (s *projectMemberStore) UpdateRole(projectID, userID, role int) error {
	result, err := s.db.Exec(`
		UPDATE project_members
		SET role = ?
		WHERE project_id = ? AND user_id = ?`,
		role,
		projectID,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotProjectMember
	}

	return nil
}

func (s *projectMemberStore) Remove(projectID, userID int) error {
	result, err := s.db.Exec(`
		DELETE FROM project_members
		WHERE project_id = ? AND user_id = ?`,
		projectID,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotProjectMember
	}

	return nil
}

func (s *projectMemberStore) TransferOwnership(projectID, currentOwnerID, newOwnerID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var role int
	err = tx.QueryRow(`
		SELECT role
		FROM project_members
		WHERE project_id = ? AND user_id = ?`,
		projectID,
		currentOwnerID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotProjectMember
	}
	if err != nil {
		return err
	}
	if role != models.RoleOwner {
		return ErrNotProjectOwner
	}

	err = tx.QueryRow(`
		SELECT role
		FROM project_members
		WHERE project_id = ? AND user_id = ?`,
		projectID,
		newOwnerID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotProjectMember
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE project_members
		SET role = ?
		WHERE project_id = ? AND user_id = ?`,
		models.RoleAdmin,
		projectID,
		currentOwnerID,
	)
	if err != nil {
		return ErrOwnershipTransferFail
	}

	_, err = tx.Exec(`
		UPDATE project_members
		SET role = ?
		WHERE project_id = ? AND user_id = ?`,
		models.RoleOwner,
		projectID,
		newOwnerID,
	)
	if err != nil {
		return ErrOwnershipTransferFail
	}

	return tx.Commit()
}
