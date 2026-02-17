package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

type Scanner interface {
	Scan(dest ...interface{}) error
}

var (
	ErrTicketNotFound   = errors.New("ticket not found")
	ErrProjectNotFound  = errors.New("project not found")
	ErrStatusNotFound   = errors.New("status not found")
	ErrPriorityNotFound = errors.New("priority not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrAssigneeNotFound = errors.New("assignee not found")
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

func validateProject(projectID int) error {
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)", projectID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate project: %w", err)
	}
	if !exists {
		return ErrProjectNotFound
	}
	return nil
}

func validateStatus(statusID int) error {
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM statuses WHERE id = ?)", statusID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate status: %w", err)
	}
	if !exists {
		return ErrStatusNotFound
	}
	return nil
}

func validatePriority(priorityID int) error {
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM priorities WHERE id = ?)", priorityID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate priority: %w", err)
	}
	if !exists {
		return ErrPriorityNotFound
	}
	return nil
}

func validateUser(userID int) error {
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate user: %w", err)
	}
	if !exists {
		return ErrUserNotFound
	}
	return nil
}

func validateAssignee(assigneeID *int) error {
	var exists bool

	if assigneeID == nil {
		return nil
	}
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", assigneeID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate asignee: %w", err)
	}
	if !exists {
		return ErrAssigneeNotFound
	}
	return nil
}

// Created an interface Scanner which is anything that has a Scan function
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

func GetAllTickets() ([]Ticket, error) {
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

func GetTicketByID(id int) (*Ticket, error) {
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

func GetTicketsByProjectID(projectID int) ([]Ticket, error) {
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

	return GetTicketByID(int(id))
}

func UpdateTicket(id int, title, description string, statusID, priorityID int, assigneeID *int) (*Ticket, error) {
	if err := validateStatus(statusID); err != nil {
		return nil, err
	}
	if err := validatePriority(priorityID); err != nil {
		return nil, err
	}
	if err := validateAssignee(assigneeID); err != nil {
		return nil, err
	}

	query := `
		UPDATE tickets
		SET title = ?, description = ?, status_id = ?, priority_id = ?, assignee_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, title, description, statusID, priorityID, assigneeID, id)
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

	return GetTicketByID(id)
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
