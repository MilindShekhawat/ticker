package models

import (
	"database/sql"
	"fmt"
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
	Position     int        `json:"position"`
	AssigneeID   *int       `json:"assignee_id"` // Pointer for nullable
	CreatedBy    int        `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"` // Pointer for nullable
}

func validateStatus(statusID int) (bool, error) {
	var statusExists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM statuses WHERE id = ?)", statusID).Scan(&statusExists)
	if err != nil {
		return false, fmt.Errorf("failed to validate status: %w", err)
	}
	return statusExists, nil
}

func GetAllTickets() ([]Ticket, error) {
	query := `
		SELECT id, project_id, ticket_number, title, description,
		       status_id, position, assignee_id, created_by,
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
		var t Ticket
		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.TicketNumber,
			&t.Title,
			&t.Description,
			&t.StatusID,
			&t.Position,
			&t.AssigneeID,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tickets: %w", err)
	}

	return tickets, nil
}

func GetTicketByID(id int) (*Ticket, error) {
	query := `
		SELECT id, project_id, ticket_number, title, description,
		       status_id, position, assignee_id, created_by,
		       created_at, updated_at, deleted_at
		FROM tickets
		WHERE id = ? AND deleted_at IS NULL
	`

	var t Ticket
	err := db.DB.QueryRow(query, id).Scan(
		&t.ID,
		&t.ProjectID,
		&t.TicketNumber,
		&t.Title,
		&t.Description,
		&t.StatusID,
		&t.Position,
		&t.AssigneeID,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return &t, nil
}

func GetTicketsByProjectID(projectID int) ([]Ticket, error) {
	query := `
		SELECT id, project_id, ticket_number, title, description,
		       status_id, position, assignee_id, created_by,
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
		var t Ticket
		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.TicketNumber,
			&t.Title,
			&t.Description,
			&t.StatusID,
			&t.Position,
			&t.AssigneeID,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tickets: %w", err)
	}

	return tickets, nil
}

func CreateTicket(projectID int, title, description string, statusID int, createdBy int, assigneeID *int) (*Ticket, error) {
	statusExists, err := validateStatus(statusID)
	if !statusExists {
		return nil, fmt.Errorf("invalid status_id: %d", statusID)
	}

	var ticketNumber int
	err = db.DB.QueryRow(`SELECT COALESCE(MAX(ticket_number), 0) + 1 WHERE project_id = ? FROM tickets`, projectID).Scan(&ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get next ticket number: %w", err)
	}

	query := `
		INSERT INTO tickets (project_id, ticket_number, title, description, status_id, created_by, assignee_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.DB.Exec(query, projectID, ticketNumber, title, description, statusID, createdBy, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return GetTicketByID(int(id))
}

func UpdateTicket(id int, title, description string, statusID int, assigneeID *int) (*Ticket, error) {
	statusExists, err := validateStatus(statusID)
	if !statusExists {
		return nil, fmt.Errorf("invalid status_id: %d", statusID)
	}

	query := `
		UPDATE tickets
		SET title = ?, description = ?, status_id = ?, assignee_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, title, description, statusID, assigneeID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("ticket with id %d not found", id)
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
		return fmt.Errorf("ticket with id %d not found", id)
	}

	return nil
}

func HardDeleteTicket(id int) error {
	query := "DELETE FROM tickets WHERE id = ?"

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("ticket with id %d not found", id)
	}

	return nil
}
