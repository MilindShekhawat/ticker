package store

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrProjectNotFound    = errors.New("project not found")
	ErrDuplicateKeyPrefix = errors.New("key prefix exists")
)

func isDuplicateError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}

func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	if !errors.As(err, &sqliteErr) {
		return false
	}
	return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
