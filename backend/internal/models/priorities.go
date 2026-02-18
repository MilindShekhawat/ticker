package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Priority struct {
	ID        int        `json:"id"`
	ProjectID *int       `json:"project_id,omitempty"`
	Label     string     `json:"label"`
	Color     string     `json:"color"`
	Position  int        `json:"position"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func scanPriority(s Scanner) (*Priority, error) {
	var p Priority
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

func ListPriorities(projectID *int) ([]Priority, error) {
	var query string
	var args []interface{}

	if projectID == nil {
		query = `
			SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
			FROM priorities
			WHERE project_id IS NULL AND deleted_at IS NULL
			ORDER BY position ASC, created_at ASC
		`
	} else {
		query = `
			SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
			FROM priorities
			WHERE project_id = ? AND deleted_at IS NULL
			ORDER BY position ASC, created_at ASC
		`
		args = append(args, *projectID)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query priorities: %w", err)
	}
	defer rows.Close()

	priorities := []Priority{}
	for rows.Next() {
		p, err := scanPriority(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan priority: %w", err)
		}
		priorities = append(priorities, *p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating priorities: %w", err)
	}

	return priorities, nil
}

func GetPriority(id int) (*Priority, error) {
	query := `
		SELECT id, project_id, label, color, position, created_at, updated_at, deleted_at
		FROM priorities
		WHERE id = ? AND deleted_at IS NULL
	`

	p, err := scanPriority(db.DB.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, ErrPriorityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get priority: %w", err)
	}

	return p, nil
}

func CreatePriority(projectID *int, label, color string) (*Priority, error) {
	if err := validateColor(color); err != nil {
		return nil, err
	}

	if projectID != nil {
		if err := validateProject(*projectID); err != nil {
			return nil, err
		}
	}

	// Get next position
	var position int
	var posQuery string
	var posArgs []interface{}

	if projectID == nil {
		posQuery = "SELECT COALESCE(MAX(position), 0) + 1 FROM priorities WHERE project_id IS NULL AND deleted_at IS NULL"
	} else {
		posQuery = "SELECT COALESCE(MAX(position), 0) + 1 FROM priorities WHERE project_id = ? AND deleted_at IS NULL"
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
		INSERT INTO priorities (project_id, label, color, position)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.DB.Exec(query, projectID, label, color, position)
	if err != nil {
		return nil, fmt.Errorf("failed to create priority: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetPriority(int(id))
}

func UpdatePriority(id int, label, color *string) (*Priority, error) {
	var updates []string
	var args []interface{}

	if label != nil {
		updates = append(updates, "label = ?")
		args = append(args, *label)
	}
	if color != nil {
		if err := validateColor(*color); err != nil {
			return nil, err
		}
		updates = append(updates, "color = ?")
		args = append(args, *color)
	}

	if len(updates) == 0 {
		return GetPriority(id)
	}

	query := fmt.Sprintf(
		"UPDATE priorities SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		strings.Join(updates, ", "),
	)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
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

	return GetPriority(id)
}

func UpdatePriorityPosition(id, position int) (*Priority, error) {
	query := `
		UPDATE priorities
		SET position = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, position, id)
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

	return GetPriority(id)
}

func DeletePriority(id int) error {
	query := `
		UPDATE priorities
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, id)
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
