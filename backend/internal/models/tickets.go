package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Ticket struct {
	ID           int        `json:"id"`
	ProjectID    int        `json:"project_id"`
	TicketNumber int        `json:"ticket_number"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	StatusID     int        `json:"status_id"`
	PriorityID   int        `json:"priority_id"`
	Position     int        `json:"position"`
	AssigneeID   *int       `json:"assignee_id"`
	CreatedBy    int        `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func scanTicket(s Scanner) (*Ticket, error) {
	var t Ticket
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

// Unused
func GetTickets() ([]Ticket, error) {
	query := `
		SELECT id, project_id, ticket_number, title, description,
		       status_id, priority_id, position, assignee_id, created_by,
		       created_at, updated_at, deleted_at
		FROM tickets
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := []Ticket{}
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, *t)
	}

	// Always check Err() after rows.Next()
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tickets: %w", err)
	}

	return tickets, nil
}

func GetTicket(id int) (*Ticket, error) {
	query := `
		SELECT id, project_id, ticket_number, title, description,
		       status_id, priority_id, position, assignee_id, created_by,
		       created_at, updated_at, deleted_at
		FROM tickets
		WHERE id = ? AND deleted_at IS NULL
	`

	t, err := scanTicket(db.DB.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return t, nil
}

func ListProjectTickets(projectID int) ([]Ticket, error) {
	query := `
		SELECT id, project_id, ticket_number, title, description,
		       status_id, priority_id, position, assignee_id, created_by,
		       created_at, updated_at, deleted_at
		FROM tickets
		WHERE project_id = ? AND deleted_at IS NULL
		ORDER BY ticket_number DESC
	`

	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := []Ticket{}
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, *t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tickets: %w", err)
	}

	return tickets, nil
}

func CreateTicket(projectID int, title, description string, statusID, priorityID, createdBy int, assigneeID *int) (*Ticket, error) {
	if err := validateProject(projectID); err != nil {
		return nil, err
	}
	if err := validateStatus(statusID); err != nil {
		return nil, err
	}
	if err := validatePriority(priorityID); err != nil {
		return nil, err
	}
	if err := validateUser(createdBy); err != nil {
		return nil, err
	}
	if err := validateAssignee(assigneeID); err != nil {
		return nil, err
	}

	var ticketNumber int
	err := db.DB.QueryRow(`SELECT COALESCE(MAX(ticket_number), 0) + 1 FROM tickets WHERE project_id = ?`, projectID).Scan(&ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get next ticket number: %w", err)
	}

	query := `
		INSERT INTO tickets (project_id, ticket_number, title, description, status_id, priority_id, created_by, assignee_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.DB.Exec(query, projectID, ticketNumber, title, description, statusID, priorityID, createdBy, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}

	return GetTicket(int(id))
}

func UpdateTicket(id int, title, description *string, statusID, priorityID, assigneeID *int) (*Ticket, error) {
	if statusID != nil {
		if err := validateStatus(*statusID); err != nil {
			return nil, err
		}
	}
	if priorityID != nil {
		if err := validatePriority(*priorityID); err != nil {
			return nil, err
		}
	}
	if err := validateAssignee(assigneeID); err != nil {
		return nil, err
	}

	var updates []string
	var args []interface{}

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

	query := fmt.Sprintf(
		"UPDATE tickets SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL",
		strings.Join(updates, ", "),
	)
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)
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

	return GetTicket(id)
}

func DeleteTicket(id int) error {
	query := `
		UPDATE tickets
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, id)
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

// Unused
func HardDeleteTicket(id int) error {
	query := "DELETE FROM tickets WHERE id = ?"

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete ticket: %w", err)
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
