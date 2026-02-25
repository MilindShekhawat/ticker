package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
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

func ListProjects() ([]Project, error) {
	rows, err := db.DB.Query(`
		SELECT id, name, description, key_prefix, created_by, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
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

func UpdateProjectByID(id int, name, description *string) (*Project, error) {
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

	query := fmt.Sprintf(`
		UPDATE projects SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		strings.Join(updates, ", "),
	)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
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
