package models

import (
	"database/sql"
	"errors"

	"github.com/MilindShekhawat/ticker/internal/db"
)

const (
	RoleOwner  = 1
	RoleAdmin  = 2
	RoleMember = 3
	RoleViewer = 4
)

var (
	ErrNotProjectMember      = errors.New("not a project member")
	ErrInvalidRole           = errors.New("invalid role")
	ErrNotProjectOwner       = errors.New("not project owner")
	ErrCannotRemoveOwner     = errors.New("cannot remove project owner")
	ErrCannotDemoteOwner     = errors.New("cannot demote project owner")
	ErrMemberAlreadyExists   = errors.New("member already exists")
	ErrOwnerAlreadyExists    = errors.New("owner already exists")
	ErrCannotPromoteToOwner  = errors.New("use transfer ownership instead")
	ErrOwnershipTransferFail = errors.New("ownership transfer failed")
)

type ProjectMember struct {
	ProjectID int `json:"project_id"`
	UserID    int `json:"user_id"`
	Role      int `json:"role"`
}

func GetUserRole(projectID, userID int) (int, error) {
	query := `
		SELECT role
		FROM project_members
		WHERE project_id = ? AND user_id = ?
	`

	var role int
	err := db.DB.QueryRow(query, projectID, userID).Scan(&role)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotProjectMember
	}
	if err != nil {
		return 0, err
	}

	return role, nil
}

func ListProjectMembers(projectID int) ([]ProjectMember, error) {
	query := `
		SELECT project_id, user_id, role
		FROM project_members
		WHERE project_id = ?
	`

	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []ProjectMember
	for rows.Next() {
		var m ProjectMember
		if err := rows.Scan(&m.ProjectID, &m.UserID, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, rows.Err()
}

func AddProjectMember(projectID, userID, role int) error {
	// Input validation
	if !IsValidRole(role) {
		return ErrInvalidRole
	}

	if role == RoleOwner {
		return ErrOwnerAlreadyExists
	}

	query := `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES (?, ?, ?)
	`

	_, err := db.DB.Exec(query, projectID, userID, role)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrMemberAlreadyExists
		}
		return err
	}

	return nil
}

func UpdateProjectMemberRole(projectID, userID, newRole int) error {
	// Input validation
	if !IsValidRole(newRole) {
		return ErrInvalidRole
	}

	if newRole == RoleOwner {
		return ErrCannotPromoteToOwner
	}

	currentRole, err := GetUserRole(projectID, userID)
	if err != nil {
		return err
	}

	if currentRole == RoleOwner {
		return ErrCannotDemoteOwner
	}

	query := `
		UPDATE project_members
		SET role = ?
		WHERE project_id = ? AND user_id = ?
	`

	_, err = db.DB.Exec(query, newRole, projectID, userID)
	return err
}

func TransferOwnership(projectID, currentOwnerID, newOwnerID int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify current owner
	var role int
	err = tx.QueryRow(
		`SELECT role FROM project_members
		 WHERE project_id = ? AND user_id = ?`,
		projectID, currentOwnerID,
	).Scan(&role)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotProjectMember
	}
	if err != nil {
		return err
	}
	if role != RoleOwner {
		return ErrNotProjectOwner
	}

	// Verify new owner is member
	err = tx.QueryRow(
		`SELECT role FROM project_members
		 WHERE project_id = ? AND user_id = ?`,
		projectID, newOwnerID,
	).Scan(&role)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotProjectMember
	}
	if err != nil {
		return err
	}

	// Demote current owner → admin
	_, err = tx.Exec(
		`UPDATE project_members
		 SET role = ?
		 WHERE project_id = ? AND user_id = ?`,
		RoleAdmin, projectID, currentOwnerID,
	)
	if err != nil {
		return ErrOwnershipTransferFail
	}

	// Promote new owner
	_, err = tx.Exec(
		`UPDATE project_members
		 SET role = ?
		 WHERE project_id = ? AND user_id = ?`,
		RoleOwner, projectID, newOwnerID,
	)
	if err != nil {
		return ErrOwnershipTransferFail
	}

	return tx.Commit()
}

func RemoveProjectMember(projectID, userID int) error {
	role, err := GetUserRole(projectID, userID)
	if err != nil {
		return err
	}

	if role == RoleOwner {
		return ErrCannotRemoveOwner
	}

	query := `
		DELETE FROM project_members
		WHERE project_id = ? AND user_id = ?
	`

	_, err = db.DB.Exec(query, projectID, userID)
	return err
}

func IsValidRole(role int) bool {
	return role >= RoleOwner && role <= RoleViewer
}
