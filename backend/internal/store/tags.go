package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

var ErrTagNotFound = errors.New("tag not found")

type TagStore interface {
	List(projectID int) ([]models.Tag, error)
	GetByID(id int) (*models.Tag, error)
	Create(projectID int, label, color string) (*models.Tag, error)
	Update(tagID int, label, color *string) (*models.Tag, error)
	Delete(tagID int) error
	GetByTicket(ticketID int) ([]models.Tag, error)
	AddToTicket(ticketID, tagID int) error
	RemoveFromTicket(ticketID, tagID int) error
	BelongsToProject(tagID, projectID int) error
}

type tagStore struct {
	db *sql.DB
}

func NewTagStore() TagStore {
	return &tagStore{db: db.DB}
}

func scanTag(s Scanner) (*models.Tag, error) {
	var t models.Tag
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

func (s *tagStore) List(projectID int) ([]models.Tag, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, label, color, created_at, updated_at, deleted_at
		FROM tags
		WHERE project_id = ? AND deleted_at IS NULL ORDER BY label ASC`,
		projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		t, err := scanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, *t)
	}
	return tags, rows.Err()
}

func (s *tagStore) GetByID(tagID int) (*models.Tag, error) {
	t, err := scanTag(s.db.QueryRow(`
		SELECT id, project_id, label, color, created_at, updated_at, deleted_at
		FROM tags
	 	WHERE id = ? AND deleted_at IS NULL`,
		tagID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTagNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return t, nil
}

func (s *tagStore) Create(projectID int, label, color string) (*models.Tag, error) {
	result, err := s.db.Exec(`
		INSERT INTO tags (project_id, label, color)
		VALUES (?, ?, ?)`,
		projectID, label, color,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return s.GetByID(int(id))
}

func (s *tagStore) Update(tagID int, label, color *string) (*models.Tag, error) {
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
		return s.GetByID(tagID)
	}

	query := fmt.Sprintf(`
		UPDATE tags SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`,
		strings.Join(updates, ", "),
	)
	args = append(args, tagID)

	result, err := s.db.Exec(query, args...)
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

	return s.GetByID(tagID)
}

func (s *tagStore) Delete(tagID int) error {
	result, err := s.db.Exec(`
		UPDATE tags SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		tagID,
	)
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

func (s *tagStore) GetByTicket(ticketID int) ([]models.Tag, error) {
	rows, err := s.db.Query(`
		SELECT t.id, t.project_id, t.label, t.color, t.created_at, t.updated_at, t.deleted_at
		FROM tags t
		JOIN ticket_tags tt ON t.id = tt.tag_id
		WHERE tt.ticket_id = ? AND t.deleted_at IS NULL
		ORDER BY t.label ASC`,
		ticketID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ticket tags: %w", err)
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		t, err := scanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, *t)
	}
	return tags, rows.Err()
}

func (s *tagStore) BelongsToProject(tagID, projectID int) error {
	var exists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM tags
			WHERE id = ? AND project_id = ? AND deleted_at IS NULL
		)`,
		tagID, projectID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTagNotFound
	}
	return nil
}

func (s *tagStore) AddToTicket(ticketID, tagID int) error {
	_, err := s.db.Exec(`
		INSERT INTO ticket_tags (ticket_id, tag_id)
		VALUES (?, ?)
		ON CONFLICT(ticket_id, tag_id) DO NOTHING`,
		ticketID, tagID,
	)
	if err != nil {
		return fmt.Errorf("failed to add tag to ticket: %w", err)
	}
	return nil
}

func (s *tagStore) RemoveFromTicket(ticketID, tagID int) error {
	result, err := s.db.Exec(`
		DELETE FROM ticket_tags WHERE ticket_id = ? AND tag_id = ?`,
		ticketID, tagID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove tag from ticket: %w", err)
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
