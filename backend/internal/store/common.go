package store

type Scanner interface {
	Scan(dest ...interface{}) error
}
