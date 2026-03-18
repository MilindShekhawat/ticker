package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentStore interface {
	ListByTicket(ticketID int) ([]models.Comment, error)
	GetByID(commentID int) (*models.Comment, error)
	GetTicketProjectID(ticketID int) (int, error)
	Create(ticketID, authorID int, body string) (*models.Comment, error)
	Update(commentID int, body string) (*models.Comment, error)
	Delete(commentID int) error
}

type commentStore struct {
	db *sql.DB
}

func NewCommentStore() CommentStore {
	return &commentStore{db: db.DB}
}

func scanComment(s Scanner) (*models.Comment, error) {
	var c models.Comment
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

func (s *commentStore) ListByTicket(ticketID int) ([]models.Comment, error) {
	rows, err := s.db.Query(`
		SELECT id, ticket_id, author_id, body, created_at, updated_at, deleted_at
		FROM comments
		WHERE ticket_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC`,
		ticketID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, *comment)
	}

	return comments, rows.Err()
}

func (s *commentStore) GetByID(commentID int) (*models.Comment, error) {
	comment, err := scanComment(s.db.QueryRow(`
		SELECT id, ticket_id, author_id, body, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = ? AND deleted_at IS NULL`,
		commentID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCommentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return comment, nil
}

func (s *commentStore) GetTicketProjectID(ticketID int) (int, error) {
	var projectID int
	err := s.db.QueryRow(`
		SELECT project_id
		FROM tickets
		WHERE id = ? AND deleted_at IS NULL`,
		ticketID,
	).Scan(&projectID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrTicketNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get ticket project: %w", err)
	}

	return projectID, nil
}

func (s *commentStore) Create(ticketID, authorID int, body string) (*models.Comment, error) {
	result, err := s.db.Exec(`
		INSERT INTO comments (ticket_id, author_id, body)
		VALUES (?, ?, ?)`,
		ticketID, authorID, body,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return s.GetByID(int(id))
}

func (s *commentStore) Update(commentID int, body string) (*models.Comment, error) {
	result, err := s.db.Exec(`
		UPDATE comments
		SET body = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		body, commentID,
	)
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

	return s.GetByID(commentID)
}

func (s *commentStore) Delete(commentID int) error {
	result, err := s.db.Exec(`
		UPDATE comments
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		commentID,
	)
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
