package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

type TicketStore interface {
	ListByProject(projectID int) ([]models.Ticket, error)
	GetByID(ticketID int) (*models.Ticket, error)
	Create(projectID int, title, description string, statusID, priorityID, createdBy int, assigneeID *int) (*models.Ticket, error)
	Update(ticketID int, title, description *string, statusID, priorityID, assigneeID *int) (*models.Ticket, error)
	Delete(ticketID int) error
	ReplaceTags(ticketID int, tagIDs []int) error
}

type ticketStore struct {
	db *sql.DB
}

func NewTicketStore() TicketStore {
	return &ticketStore{db: db.DB}
}

func scanTicket(s Scanner) (*models.Ticket, error) {
	var t models.Ticket
	err := s.Scan(
		&t.ID,
		&t.ProjectID,
		&t.TicketNumber,
		&t.Title,
		&t.Description,
		&t.StatusID,
		&t.PriorityID,
		&t.Position,
		&t.AssigneeID,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
	)
	return &t, err
}

var (
	ErrTicketNotFound = errors.New("ticket not found")
)

func (s *ticketStore) ListByProject(projectID int) ([]models.Ticket, error) {
	rows, err := s.db.Query(`
		SELECT id, project_id, ticket_number, title, description,
		       status_id, priority_id, position, assignee_id, created_by,
		       created_at, updated_at, deleted_at
		FROM tickets
		WHERE project_id = ? AND deleted_at IS NULL
		ORDER BY ticket_number DESC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]models.Ticket, 0)
	for rows.Next() {
		ticket, err := scanTicket(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, *ticket)
	}

	return tickets, rows.Err()
}

func (s *ticketStore) GetByID(ticketID int) (*models.Ticket, error) {
	ticket, err := scanTicket(s.db.QueryRow(`
		SELECT id, project_id, ticket_number, title, description,
		       status_id, priority_id, position, assignee_id, created_by,
		       created_at, updated_at, deleted_at
		FROM tickets
	 	WHERE id = ? AND deleted_at IS NULL`,
		ticketID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}

func (s *ticketStore) Create(projectID int, title, description string, statusID, priorityID, createdBy int, assigneeID *int) (*models.Ticket, error) {
	var ticketNumber int
	err := s.db.QueryRow(`
		SELECT COALESCE(MAX(ticket_number), 0) + 1
		FROM tickets
		WHERE project_id = ?`,
		projectID,
	).Scan(&ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get next ticket number: %w", err)
	}

	result, err := s.db.Exec(`
		INSERT INTO tickets (project_id, ticket_number, title, description, status_id, priority_id, created_by, assignee_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		projectID, ticketNumber, title, description, statusID, priorityID, createdBy, assigneeID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return s.GetByID(int(id))
}

func (s *ticketStore) Update(ticketID int, title, description *string, statusID, priorityID, assigneeID *int) (*models.Ticket, error) {
	var updates []string
	var args []any

	if title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *title)
	}
	if description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *description)
	}
	if statusID != nil {
		updates = append(updates, "status_id = ?")
		args = append(args, *statusID)
	}
	if priorityID != nil {
		updates = append(updates, "priority_id = ?")
		args = append(args, *priorityID)
	}
	if assigneeID != nil {
		updates = append(updates, "assignee_id = ?")
		args = append(args, *assigneeID)
	}

	if len(updates) == 0 {
		return s.GetByID(ticketID)
	}

	query := fmt.Sprintf(
		"UPDATE tickets SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		strings.Join(updates, ", "),
	)
	args = append(args, ticketID)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrTicketNotFound
	}

	return s.GetByID(ticketID)
}

func (s *ticketStore) Delete(ticketID int) error {
	result, err := s.db.Exec(`
		UPDATE tickets SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`,
		ticketID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrTicketNotFound
	}

	return nil
}

func (s *ticketStore) ReplaceTags(ticketID int, tagIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		DELETE FROM ticket_tags
		WHERE ticket_id = ?`,
		ticketID,
	); err != nil {
		return fmt.Errorf("failed to clear ticket tags: %w", err)
	}

	for _, tagID := range tagIDs {
		if _, err := tx.Exec(`
			INSERT INTO ticket_tags (ticket_id, tag_id)
			VALUES (?, ?)
			ON CONFLICT(ticket_id, tag_id) DO NOTHING`,
			ticketID, tagID,
		); err != nil {
			return fmt.Errorf("failed to insert ticket tag: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
