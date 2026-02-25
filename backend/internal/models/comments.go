package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Comment struct {
	ID        int        `json:"id"`
	TicketID  int        `json:"ticket_id"`
	AuthorID  int        `json:"author_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func scanComment(s Scanner) (*Comment, error) {
	var c Comment
	err := s.Scan(
		&c.ID,
		&c.TicketID,
		&c.AuthorID,
		&c.Body,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	return &c, err
}

func ListComments(ticketID int) ([]Comment, error) {
	query := `
		SELECT id, ticket_id, author_id, body, created_at, updated_at, deleted_at
		FROM comments
		WHERE ticket_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := db.DB.Query(query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, *c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

func GetComment(id int) (*Comment, error) {
	query := `
		SELECT id, ticket_id, author_id, body, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = ? AND deleted_at IS NULL
	`

	c, err := scanComment(db.DB.QueryRow(query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCommentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return c, nil
}

func CreateComment(ticketID, authorID int, body string) (*Comment, error) {
	if strings.TrimSpace(body) == "" {
		return nil, ErrInvalidCommentBody
	}

	if err := doesTicketExists(ticketID); err != nil {
		return nil, err
	}
	if err := doesUserExists(authorID); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO comments (ticket_id, author_id, body)
		VALUES (?, ?, ?)
	`

	result, err := db.DB.Exec(query, ticketID, authorID, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetComment(int(id))
}

func UpdateComment(id int, body string) (*Comment, error) {
	query := `
		UPDATE comments
		SET body = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, body, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrCommentNotFound
	}

	return GetComment(id)
}

func DeleteComment(id int) error {
	query := `
		UPDATE comments
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrCommentNotFound
	}

	return nil
}
