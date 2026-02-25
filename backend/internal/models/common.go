package models

import (
	"errors"
	"regexp"

	"github.com/MilindShekhawat/ticker/internal/db"
)

// Common errors used across models
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrAssigneeNotFound   = errors.New("assignee not found")
	ErrProjectNotFound    = errors.New("project not found")
	ErrTicketNotFound     = errors.New("ticket not found")
	ErrPriorityNotFound   = errors.New("priority not found")
	ErrStatusNotFound     = errors.New("status not found")
	ErrTagNotFound        = errors.New("tag not found")
	ErrDuplicateKeyPrefix = errors.New("key prefix already exists")
	ErrUserPrefsNotFound  = errors.New("user preferences not found")
	ErrCommentNotFound    = errors.New("comment not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCommentBody = errors.New("invalid comment body")
	ErrInvalidColor       = errors.New("invalid color")
)

// Scanner is an interface for anything that can Scan (sql.Row, sql.Rows, etc.)
type Scanner interface {
	Scan(dest ...interface{}) error
}

func doesUserExists(userID int) error {
	var exists bool

	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		userID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}
	return nil
}

func doesAssigneeExists(assigneeID *int) error {
	if assigneeID == nil {
		return nil
	}

	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		*assigneeID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAssigneeNotFound
	}
	return nil
}

func doesProjectExists(projectID int) error {
	var exists bool

	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)",
		projectID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrProjectNotFound
	}
	return nil
}

func doesKeyPrefixExists(keyPrefix string) error {
	var exists bool

	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM projects WHERE key_prefix = ?)",
		keyPrefix,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateKeyPrefix
	}
	return nil
}

func doesTicketExists(ticketID int) error {
	var exists bool

	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM tickets WHERE id = ? AND deleted_at IS NULL)",
		ticketID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTicketNotFound
	}
	return nil
}

func doesStatusExists(statusID int) error {
	var exists bool

	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM statuses WHERE id = ? AND deleted_at IS NULL)",
		statusID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrStatusNotFound
	}
	return nil
}

func doesPriorityExists(priorityID int) error {
	var exists bool

	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM priorities WHERE id = ? AND deleted_at IS NULL)",
		priorityID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPriorityNotFound
	}
	return nil
}

func isValidColor(color string) error {
	if color == "" {
		return nil
	}

	matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, color)
	if !matched {
		return ErrInvalidColor
	}

	return nil
}
