package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Status struct {
	ID        int        `json:"id"`
	ProjectID *int       `json:"project_id,omitempty"`
	Label     string     `json:"label"`
	Color     string     `json:"color"`
	Position  int        `json:"position"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func scanStatus(s Scanner) (*Status, error) {
	var st Status
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

func ListStatuses(projectID *int) ([]Status, error) {
	var query string
	var args []interface{}

	if projectID == nil {
		query = `
			SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
			FROM statuses
			WHERE project_id IS NULL AND deleted_at IS NULL
			ORDER BY position ASC, created_at ASC
		`
	} else {
		query = `
			SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
			FROM statuses
			WHERE project_id = ? AND deleted_at IS NULL
			ORDER BY position ASC, created_at ASC
		`
		args = append(args, *projectID)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query statuses: %w", err)
	}
	defer rows.Close()

	statuses := []Status{}
	for rows.Next() {
		st, err := scanStatus(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status: %w", err)
		}
		statuses = append(statuses, *st)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating statuses: %w", err)
	}

	return statuses, nil
}

func GetStatus(id int) (*Status, error) {
	query := `
		SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
		FROM statuses
		WHERE id = ? AND deleted_at IS NULL
	`

	st, err := scanStatus(db.DB.QueryRow(query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStatusNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return st, nil
}

func CreateStatus(projectID *int, label, color string) (*Status, error) {
	if err := isValidColor(color); err != nil {
		return nil, err
	}

	if projectID != nil {
		if err := doesProjectExists(*projectID); err != nil {
			return nil, err
		}
	}

	// Get next position
	var position int
	var posQuery string
	var posArgs []interface{}

	if projectID == nil {
		posQuery = "SELECT COALESCE(MAX(position), 0) + 1 FROM statuses WHERE project_id IS NULL AND deleted_at IS NULL"
	} else {
		posQuery = "SELECT COALESCE(MAX(position), 0) + 1 FROM statuses WHERE project_id = ? AND deleted_at IS NULL"
		posArgs = append(posArgs, *projectID)
	}

	err := db.DB.QueryRow(posQuery, posArgs...).Scan(&position)
	if err != nil {
		return nil, fmt.Errorf("failed to get next position: %w", err)
	}

	// Set default color if empty
	if color == "" {
		color = "#808080"
	}

	query := `
		INSERT INTO statuses (project_id, label, color, position)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.DB.Exec(query, projectID, label, color, position)
	if err != nil {
		return nil, fmt.Errorf("failed to create status: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetStatus(int(id))
}

func UpdateStatus(id int, label, color *string) (*Status, error) {
	var updates []string
	var args []interface{}

	if label != nil {
		updates = append(updates, "label = ?")
		args = append(args, *label)
	}
	if color != nil {
		if err := isValidColor(*color); err != nil {
			return nil, err
		}
		updates = append(updates, "color = ?")
		args = append(args, *color)
	}

	if len(updates) == 0 {
		return GetStatus(id)
	}

	query := fmt.Sprintf(
		"UPDATE statuses SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		strings.Join(updates, ", "),
	)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
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

	return GetStatus(id)
}

func UpdateStatusPosition(id, position int) (*Status, error) {
	query := `
		UPDATE statuses
		SET position = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, position, id)
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

	return GetStatus(id)
}

func DeleteStatus(id int) error {
	query := `
		UPDATE statuses
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, id)
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
