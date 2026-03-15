package services

import "errors"

var (
	ErrInvalidProject       = errors.New("invalid project")
	ErrInvalidProjectName   = errors.New("invalid project name")
	ErrInvalidKeyPrefix     = errors.New("invalid key prefix")
	ErrInvalidProjectUpdate = errors.New("no fields to update")
)
