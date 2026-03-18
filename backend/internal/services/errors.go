package services

import "errors"

var (
	ErrInvalidProject       = errors.New("invalid project")
	ErrInvalidProjectName   = errors.New("invalid project name")
	ErrInvalidKeyPrefix     = errors.New("invalid key prefix")
	ErrInvalidProjectUpdate = errors.New("no fields to update")
	ErrInvalidCommentBody   = errors.New("invalid comment body")
	ErrTicketTitleRequired  = errors.New("title is required")
	ErrTicketTitleTooLong   = errors.New("title too long")
	ErrInvalidStatusID      = errors.New("invalid status ID")
	ErrInvalidPriorityID    = errors.New("invalid priority ID")
	ErrInvalidAssigneeID    = errors.New("invalid assignee ID")
)
