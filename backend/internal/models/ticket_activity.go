package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MilindShekhawat/ticker/internal/db"
)

// Activity types
const (
	ActivityCreated            = 1
	ActivityStatusChanged      = 2
	ActivityAssigneeChanged    = 3
	ActivityCommented          = 4
	ActivityTitleChanged       = 5
	ActivityDescriptionChanged = 6
	ActivityTagAdded           = 7
	ActivityTagRemoved         = 8
	ActivityPriorityChanged    = 9
)

type TicketActivity struct {
	ID           int                    `json:"id"`
	TicketID     int                    `json:"ticket_id"`
	ActorID      int                    `json:"actor_id"`
	ActivityType int                    `json:"activity_type"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

func scanActivity(s Scanner) (*TicketActivity, error) {
	var a TicketActivity
	var metadataJSON []byte

	err := s.Scan(
		&a.ID,
		&a.TicketID,
		&a.ActorID,
		&a.ActivityType,
		&metadataJSON,
		&a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse JSON metadata if present
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &a.Metadata); err != nil {
			return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
		}
	}

	return &a, nil
}

func ListActivityByTicket(ticketID int) ([]TicketActivity, error) {
	query := `
		SELECT id, ticket_id, actor_id, activity_type, metadata, created_at
		FROM ticket_activity
		WHERE ticket_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity: %w", err)
	}
	defer rows.Close()

	activities := []TicketActivity{}
	for rows.Next() {
		a, err := scanActivity(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, *a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activity: %w", err)
	}

	return activities, nil
}

func CreateActivity(ticketID, actorID, activityType int, metadata map[string]interface{}) error {
	// Validate ticket and actor
	if err := doesTicketExists(ticketID); err != nil {
		return err
	}
	if err := doesUserExists(actorID); err != nil {
		return err
	}

	// Convert metadata to JSON
	var metadataJSON []byte
	var err error
	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	query := `
		INSERT INTO ticket_activity (ticket_id, actor_id, activity_type, metadata)
		VALUES (?, ?, ?, ?)
	`

	_, err = db.DB.Exec(query, ticketID, actorID, activityType, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

// Helper functions for common activity types

func LogTicketCreated(ticketID, actorID int) error {
	return CreateActivity(ticketID, actorID, ActivityCreated, nil)
}

func LogStatusChanged(ticketID, actorID, fromStatusID, toStatusID int) error {
	return CreateActivity(ticketID, actorID, ActivityStatusChanged, map[string]interface{}{
		"from_status_id": fromStatusID,
		"to_status_id":   toStatusID,
	})
}

func LogAssigneeChanged(ticketID, actorID int, fromAssigneeID, toAssigneeID *int) error {
	metadata := map[string]interface{}{}
	if fromAssigneeID != nil {
		metadata["from_assignee_id"] = *fromAssigneeID
	}
	if toAssigneeID != nil {
		metadata["to_assignee_id"] = *toAssigneeID
	}
	return CreateActivity(ticketID, actorID, ActivityAssigneeChanged, metadata)
}

func LogCommented(ticketID, actorID, commentID int) error {
	return CreateActivity(ticketID, actorID, ActivityCommented, map[string]interface{}{
		"comment_id": commentID,
	})
}

func LogTitleChanged(ticketID, actorID int, oldTitle, newTitle string) error {
	return CreateActivity(ticketID, actorID, ActivityTitleChanged, map[string]interface{}{
		"old": oldTitle,
		"new": newTitle,
	})
}

func LogDescriptionChanged(ticketID, actorID int, oldDesc, newDesc string) error {
	return CreateActivity(ticketID, actorID, ActivityDescriptionChanged, map[string]interface{}{
		"old": oldDesc,
		"new": newDesc,
	})
}

func LogTagAdded(ticketID, actorID, tagID int) error {
	return CreateActivity(ticketID, actorID, ActivityTagAdded, map[string]interface{}{
		"tag_id": tagID,
	})
}

func LogTagRemoved(ticketID, actorID, tagID int) error {
	return CreateActivity(ticketID, actorID, ActivityTagRemoved, map[string]interface{}{
		"tag_id": tagID,
	})
}

func LogPriorityChanged(ticketID, actorID int, fromPriorityID, toPriorityID *int) error {
	metadata := map[string]interface{}{}
	if fromPriorityID != nil {
		metadata["from_priority_id"] = *fromPriorityID
	}
	if toPriorityID != nil {
		metadata["to_priority_id"] = *toPriorityID
	}
	return CreateActivity(ticketID, actorID, ActivityPriorityChanged, metadata)
}
