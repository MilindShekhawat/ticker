package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

var ErrStatusNotFound = errors.New("status not found")

type StatusStore interface {
	List(projectID int) ([]models.Status, error)
	GetByID(id int) (*models.Status, error)
	Create(projectID int, label, color string) (*models.Status, error)
	Update(id int, label, color *string) (*models.Status, error)
	UpdatePosition(id, position int) (*models.Status, error)
	Delete(id int) error
}

type statusStore struct {
	db *sql.DB
}

func NewStatusStore() StatusStore {
	return &statusStore{db: db.DB}
}

func scanStatus(s Scanner) (*models.Status, error) {
	var st models.Status
	err := s.Scan(
		&st.ID,
		&st.ProjectID,
		&st.Label,
		&st.Color,
		&st.Position,
		&st.CreatedAt,
		&st.UpdatedAt,
		&st.DeletedAt,
	)
	return &st, err
}

func (s *statusStore) List(projectID int) ([]models.Status, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
		FROM statuses
		WHERE project_id = ? AND deleted_at IS NULL ORDER BY position ASC, created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query statuses: %w", err)
	}
	defer rows.Close()

	statuses := []models.Status{}
	for rows.Next() {
		st, err := scanStatus(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status: %w", err)
		}
		statuses = append(statuses, *st)
	}
	return statuses, rows.Err()
}

func (s *statusStore) GetByID(statusID int) (*models.Status, error) {
	st, err := scanStatus(s.db.QueryRow(`
		SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
		FROM statuses
		WHERE id = ? AND deleted_at IS NULL`,
		statusID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStatusNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}
	return st, nil
}

func (s *statusStore) Create(projectID int, label, color string) (*models.Status, error) {
	var position int
	err := s.db.QueryRow(
		`SELECT COALESCE(MAX(position), 0) + 1 FROM statuses WHERE project_id = ? AND deleted_at IS NULL`,
		projectID,
	).Scan(&position)
	if err != nil {
		return nil, fmt.Errorf("failed to get next position: %w", err)
	}

	result, err := s.db.Exec(`
		INSERT INTO statuses (project_id, label, color, position)
		VALUES (?, ?, ?, ?)`,
		projectID, label, color, position,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create status: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return s.GetByID(int(id))
}

func (s *statusStore) Update(statusID int, label, color *string) (*models.Status, error) {
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
		return s.GetByID(statusID)
	}

	query := fmt.Sprintf(
		`UPDATE statuses SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`,
		strings.Join(updates, ", "),
	)
	args = append(args, statusID)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrStatusNotFound
	}

	return s.GetByID(statusID)
}

func (s *statusStore) UpdatePosition(statusID, position int) (*models.Status, error) {
	result, err := s.db.Exec(`
		UPDATE statuses SET position = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		position, statusID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update status position: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrStatusNotFound
	}

	return s.GetByID(statusID)
}

func (s *statusStore) Delete(statusID int) error {
	result, err := s.db.Exec(`
		UPDATE statuses SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`, statusID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrStatusNotFound
	}

	return nil
}
