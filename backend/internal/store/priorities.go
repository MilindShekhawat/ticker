package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

var ErrPriorityNotFound = errors.New("priority not found")

type PriorityStore interface {
	List(projectID int) ([]models.Priority, error)
	GetByID(priorityID int) (*models.Priority, error)
	Create(projectID int, label, color string) (*models.Priority, error)
	Update(priorityID int, label, color *string) (*models.Priority, error)
	UpdatePosition(priorityID, position int) (*models.Priority, error)
	Delete(priorityID int) error
}

type priorityStore struct {
	db *sql.DB
}

func NewPriorityStore() PriorityStore {
	return &priorityStore{db: db.DB}
}

func scanPriority(s Scanner) (*models.Priority, error) {
	var p models.Priority
	err := s.Scan(
		&p.ID,
		&p.ProjectID,
		&p.Label,
		&p.Color,
		&p.Position,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	return &p, err
}

func (s *priorityStore) List(projectID int) ([]models.Priority, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
		FROM priorities
		WHERE project_id = ? AND deleted_at IS NULL
		ORDER BY position ASC, created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query priorities: %w", err)
	}
	defer rows.Close()

	priorities := []models.Priority{}
	for rows.Next() {
		p, err := scanPriority(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan priority: %w", err)
		}
		priorities = append(priorities, *p)
	}
	return priorities, rows.Err()
}

func (s *priorityStore) GetByID(priorityID int) (*models.Priority, error) {
	p, err := scanPriority(s.db.QueryRow(`
		SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
		FROM priorities
		WHERE id = ? AND deleted_at IS NULL`,
		priorityID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPriorityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get priority: %w", err)
	}
	return p, nil
}

func (s *priorityStore) Create(projectID int, label, color string) (*models.Priority, error) {
	var position int
	err := s.db.QueryRow(`
		SELECT COALESCE(MAX(position), 0) + 1
		FROM priorities
		WHERE project_id = ? AND deleted_at IS NULL`,
		projectID,
	).Scan(&position)
	if err != nil {
		return nil, fmt.Errorf("failed to get next position: %w", err)
	}

	result, err := s.db.Exec(`
		INSERT INTO priorities (project_id, label, color, position)
		VALUES (?, ?, ?, ?)`,
		projectID, label, color, position,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create priority: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return s.GetByID(int(id))
}

func (s *priorityStore) Update(priorityID int, label, color *string) (*models.Priority, error) {
	var updates []string
	var args []interface{}

	if label != nil {
		updates = append(updates, "label = ?")
		args = append(args, *label)
	}
	if color != nil {
		updates = append(updates, "color = ?")
		args = append(args, *color)
	}

	if len(updates) == 0 {
		return s.GetByID(priorityID)
	}

	query := fmt.Sprintf(
		`UPDATE priorities SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`,
		strings.Join(updates, ", "),
	)
	args = append(args, priorityID)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update priority: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrPriorityNotFound
	}

	return s.GetByID(priorityID)
}

func (s *priorityStore) UpdatePosition(priorityID, position int) (*models.Priority, error) {
	result, err := s.db.Exec(`
		UPDATE priorities
		SET position = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		position, priorityID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update priority position: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrPriorityNotFound
	}

	return s.GetByID(priorityID)
}

func (s *priorityStore) Delete(priorityID int) error {
	result, err := s.db.Exec(`
		UPDATE priorities
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		priorityID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete priority: %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrPriorityNotFound
	}

	return nil
}
