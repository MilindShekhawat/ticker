package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/MilindShekhawat/ticker/internal/db"
)

// Common errors used across models
var (
	ErrTicketNotFound     = errors.New("ticket not found")
	ErrProjectNotFound    = errors.New("project not found")
	ErrStatusNotFound     = errors.New("status not found")
	ErrPriorityNotFound   = errors.New("priority not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrAssigneeNotFound   = errors.New("assignee not found")
	ErrDuplicateKeyPrefix = errors.New("key prefix already exists")
	ErrTagNotFound        = errors.New("tag not found")
	ErrUserPrefsNotFound  = errors.New("user preferences not found")
	ErrCommentNotFound    = errors.New("comment not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

// Scanner is an interface for anything that can Scan (sql.Row, sql.Rows, etc.)
type Scanner interface {
	Scan(dest ...interface{}) error
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
	if assigneeID == nil {
		return nil
	}
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", assigneeID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate assignee: %w", err)
	}
	if !exists {
		return ErrAssigneeNotFound
	}
	return nil
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

func validateTicket(ticketID int) error {
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tickets WHERE id = ? AND deleted_at IS NULL)", ticketID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate ticket: %w", err)
	}
	if !exists {
		return ErrTicketNotFound
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

func ValidateKeyPrefix(keyPrefix string) error {
	var exists bool

	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE key_prefix = ? AND deleted_at IS NULL)", keyPrefix).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check key prefix: %w", err)
	}
	if exists {
		return ErrDuplicateKeyPrefix
	}
	return nil
}

func validateColor(color string) error {
	if color == "" {
		return nil // Allow empty, will use default
	}
	matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, color)
	if !matched {
		return fmt.Errorf("invalid color format: must be #RRGGBB")
	}
	return nil
}
