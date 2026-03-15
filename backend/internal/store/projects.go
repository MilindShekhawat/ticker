package store

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

type ProjectStore interface {
	ListByUser(userID int) ([]models.Project, error)
	GetByUser(projectID, userID int) (*models.Project, error)
	GetByID(id int) (*models.Project, error)
	Create(name, description, keyPrefix string, userId int) (*models.Project, error)
	UpdateByUser(id, userID int, name, description *string) (*models.Project, error)
	DeleteByUser(id, userID int) error
}

type projectStore struct {
	db *sql.DB
}

func NewProjectStore() ProjectStore {
	return &projectStore{db: db.DB}
}

func scanProject(s Scanner) (*models.Project, error) {
	var p models.Project
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

func (s *projectStore) ListByUser(userID int) ([]models.Project, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.name, p.description, p.key_prefix, p.created_by, p.created_at, p.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = ?
		ORDER BY p.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}

	return projects, rows.Err()
}

func (s *projectStore) GetByUser(projectID, userID int) (*models.Project, error) {
	p, err := scanProject(s.db.QueryRow(`
		SELECT p.id, p.name, p.description, p.key_prefix, p.created_by, p.created_at, p.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE p.id = ? AND pm.user_id = ?`,
		projectID, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *projectStore) GetByID(id int) (*models.Project, error) {
	p, err := scanProject(s.db.QueryRow(`
		SELECT id, name, description, key_prefix, created_by, created_at, updated_at
		FROM projects WHERE id = ?`,
		id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *projectStore) Create(name, description, keyPrefix string, createdBy int) (*models.Project, error) {
	tx, err := s.db.Begin()
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
		if isDuplicateError(err) {
			return nil, ErrDuplicateKeyPrefix
		}
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
		projectID, createdBy, models.RoleOwner,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetByID(int(projectID))
}

func (s *projectStore) UpdateByUser(id, userID int, name, description *string) (*models.Project, error) {
	var updates []string
	var args []any

	if name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *name)
	}
	if description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *description)
	}

	args = append(args, id, userID)

	result, err := s.db.Exec(`
		UPDATE projects
		SET `+strings.Join(updates, ", ")+`, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND EXISTS (
			SELECT 1 FROM project_members pm
			WHERE pm.project_id = projects.id AND pm.user_id = ?
		  )`,
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

	return s.GetByUser(id, userID)
}

func (s *projectStore) DeleteByUser(id, userID int) error {
	result, err := s.db.Exec(`
		DELETE FROM projects
		WHERE id = ?
		  AND EXISTS (
			SELECT 1 FROM project_members pm
			WHERE pm.project_id = projects.id AND pm.user_id = ?
		  )`,
		id, userID,
	)
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
