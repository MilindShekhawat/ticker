package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/MilindShekhawat/ticker/internal/db"
)

var (
	ErrInvalidProject       = errors.New("invalid project")
	ErrInvalidProjectName   = errors.New("invalid project name")
	ErrInvalidKeyPrefix     = errors.New("invalid key prefix")
	ErrInvalidProjectUpdate = errors.New("no fields to update")
)

const (
	MaxProjectNameLength = 30
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	KeyPrefix   string    `json:"key_prefix"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	)
	return &p, err
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

func ListProjects() ([]Project, error) {
	rows, err := db.DB.Query(`
		SELECT id, name, description, key_prefix, created_by, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func GetProjectByID(id int) (*Project, error) {
	p, err := scanProject(db.DB.QueryRow(`
		SELECT id, name, description, key_prefix, created_by, created_at, updated_at
		FROM projects
		WHERE id = ?`,
		id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func CreateProject(name, description, keyPrefix string, createdBy int) (*Project, error) {
	name = strings.TrimSpace(name)
	keyPrefix = strings.ToUpper(strings.TrimSpace(keyPrefix))

	if name == "" || len(name) > MaxProjectNameLength {
		return nil, ErrInvalidProjectName
	}
	if !isValidKeyPrefix(keyPrefix) {
		return nil, ErrInvalidKeyPrefix
	}
	if createdBy == 0 {
		return nil, ErrInvalidProject
	}

	if err := doesUserExists(createdBy); err != nil {
		return nil, err
	}
	if err := doesKeyPrefixExists(keyPrefix); err != nil {
		return nil, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(`
		INSERT INTO projects (name, description, key_prefix, created_by)
		VALUES (?, ?, ?, ?)`,
		name, description, keyPrefix, createdBy,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO project_members (project_id, user_id, role)
		VALUES (?, ?, ?)`,
		projectID, createdBy, RoleOwner,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return GetProjectByID(int(projectID))
}

func UpdateProject(id int, name, description *string) (*Project, error) {
	var updates []string
	var args []any

	if name != nil {
		n := strings.TrimSpace(*name)
		if n == "" || len(n) > MaxProjectNameLength {
			return nil, ErrInvalidProjectName
		}
		updates = append(updates, "name = ?")
		args = append(args, n)
	}

	if description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *description)
	}

	if len(updates) == 0 {
		return nil, ErrInvalidProjectUpdate
	}

	args = append(args, id)

	result, err := db.DB.Exec(`
		UPDATE projects
		SET `+strings.Join(updates, ", ")+`, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		args...)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, ErrProjectNotFound
	}

	return GetProjectByID(id)
}

func DeleteProjectByID(id int) error {
	result, err := db.DB.Exec(`DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrProjectNotFound
	}

	return nil
}
