package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Project struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	KeyPrefix   string     `json:"key_prefix"`
	CreatedBy   int        `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func scanProject(s Scanner) (*Project, error) {
	var p Project
	err := s.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.KeyPrefix,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	return &p, err
}

func ListProjects() ([]Project, error) {
	query := `
		SELECT id, name, description, key_prefix, created_by,
		       created_at, updated_at, deleted_at
		FROM projects
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	projects := []Project{}
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, *p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

func GetProjectByID(id int) (*Project, error) {
	query := `
		SELECT id, name, description, key_prefix, created_by,
		       created_at, updated_at, deleted_at
		FROM projects
		WHERE id = ? AND deleted_at IS NULL
	`

	p, err := scanProject(db.DB.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return p, nil
}

func CreateProject(name, description, keyPrefix string, createdBy int) (*Project, error) {
	if err := validateUser(createdBy); err != nil {
		return nil, err
	}
	if err := ValidateKeyPrefix(keyPrefix); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO projects (name, description, key_prefix, created_by)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.DB.Exec(query, name, description, keyPrefix, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetProjectByID(int(id))
}

func UpdateProject(id int, name, description *string) (*Project, error) {
	var updates []string
	var args []interface{}

	if name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *name)
	}
	if description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *description)
	}

	query := fmt.Sprintf(
		"UPDATE projects SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		strings.Join(updates, ", "),
	)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrProjectNotFound
	}

	return GetProjectByID(id)
}

func DeleteProject(id int) error {
	query := `
		UPDATE projects
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrProjectNotFound
	}

	return nil
}
