package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Tag struct {
	ID        int        `json:"id"`
	ProjectID *int       `json:"project_id,omitempty"`
	Label     string     `json:"label"`
	Color     string     `json:"color"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func scanTag(s Scanner) (*Tag, error) {
	var t Tag
	err := s.Scan(
		&t.ID,
		&t.ProjectID,
		&t.Label,
		&t.Color,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
	)
	return &t, err
}

func ListTags(projectID *int) ([]Tag, error) {
	var query string
	var args []interface{}

	if projectID == nil {
		query = `
			SELECT id, project_id, label, color, created_at, updated_at, deleted_at
			FROM tags
			WHERE project_id IS NULL AND deleted_at IS NULL
			ORDER BY label ASC
		`
	} else {
		query = `
			SELECT id, project_id, label, color, created_at, updated_at, deleted_at
			FROM tags
			WHERE project_id = ? AND deleted_at IS NULL
			ORDER BY label ASC
		`
		args = append(args, *projectID)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		t, err := scanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, *t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tags: %w", err)
	}

	return tags, nil
}

func GetTag(id int) (*Tag, error) {
	query := `
		SELECT id, project_id, label, color, created_at, updated_at, deleted_at
		FROM tags
		WHERE id = ? AND deleted_at IS NULL
	`

	t, err := scanTag(db.DB.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return t, nil
}

func CreateTag(projectID *int, label, color string) (*Tag, error) {
	if err := validateColor(color); err != nil {
		return nil, err
	}

	if projectID != nil {
		if err := validateProject(*projectID); err != nil {
			return nil, err
		}
	}

	// Set default color if empty
	if color == "" {
		color = "#808080"
	}

	query := `
		INSERT INTO tags (project_id, label, color)
		VALUES (?, ?, ?)
	`
	result, err := db.DB.Exec(query, projectID, label, color)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetTag(int(id))
}

func UpdateTag(id int, label, color *string) (*Tag, error) {
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
		return GetTag(id)
	}

	query := fmt.Sprintf(
		"UPDATE tags SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		strings.Join(updates, ", "),
	)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrTagNotFound
	}

	return GetTag(id)
}

func DeleteTag(id int) error {
	query := `
		UPDATE tags
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrTagNotFound
	}

	return nil
}

// Ticket Tags Operations

func GetTagsByTicketID(ticketID int) ([]Tag, error) {
	query := `
		SELECT t.id, t.project_id, t.label, t.color, t.created_at, t.updated_at, t.deleted_at
		FROM tags t
		JOIN ticket_tags tt ON t.id = tt.tag_id
		WHERE tt.ticket_id = ? AND t.deleted_at IS NULL
		ORDER BY t.label ASC
	`

	rows, err := db.DB.Query(query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ticket tags: %w", err)
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		t, err := scanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, *t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tags: %w", err)
	}

	return tags, nil
}

func AddTagToTicket(ticketID, tagID int) error {
	// Verify tag exists and is not deleted
	tag, err := GetTag(tagID)
	if err != nil {
		return err
	}
	if tag.DeletedAt != nil {
		return fmt.Errorf("cannot add deleted tag to ticket")
	}

	query := `
		INSERT INTO ticket_tags (ticket_id, tag_id)
		VALUES (?, ?)
		ON CONFLICT(ticket_id, tag_id) DO NOTHING
	`

	_, err = db.DB.Exec(query, ticketID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to ticket: %w", err)
	}

	return nil
}

func RemoveTagFromTicket(ticketID, tagID int) error {
	query := `DELETE FROM ticket_tags WHERE ticket_id = ? AND tag_id = ?`

	result, err := db.DB.Exec(query, ticketID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from ticket: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tag not found on ticket")
	}

	return nil
}
