package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MilindShekhawat/ticker/internal/db"
	"github.com/MilindShekhawat/ticker/internal/models"
)

type TicketActivityStore interface {
	ListByTicket(ticketID int) ([]models.TicketActivity, error)
	GetTicketProjectID(ticketID int) (int, error)
}

type ticketActivityStore struct {
	db *sql.DB
}

func NewTicketActivityStore() TicketActivityStore {
	return &ticketActivityStore{db: db.DB}
}

func scanTicketActivity(s Scanner) (*models.TicketActivity, error) {
	var activity models.TicketActivity
	var metadataJSON []byte

	err := s.Scan(
		&activity.ID,
		&activity.TicketID,
		&activity.ActorID,
		&activity.ActivityType,
		&metadataJSON,
		&activity.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &activity.Metadata); err != nil {
			return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
		}
	}

	return &activity, nil
}

func (s *ticketActivityStore) ListByTicket(ticketID int) ([]models.TicketActivity, error) {
	rows, err := s.db.Query(`
		SELECT id, ticket_id, actor_id, activity_type, metadata, created_at
		FROM ticket_activity
		WHERE ticket_id = ?
		ORDER BY created_at DESC`,
		ticketID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity: %w", err)
	}
	defer rows.Close()

	activities := make([]models.TicketActivity, 0)
	for rows.Next() {
		activity, err := scanTicketActivity(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, *activity)
	}

	return activities, rows.Err()
}

func (s *ticketActivityStore) GetTicketProjectID(ticketID int) (int, error) {
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
