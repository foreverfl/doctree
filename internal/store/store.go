package store

import _ "embed"

//go:embed schema.sql
var schemaSQL string

// Store wraps the SQLite handle used by the daemon.
// All writes go through BEGIN IMMEDIATE so port allocation is race-free.
type Store struct {
	// db *sql.DB
}

// Open opens (and migrates) the SQLite database at path.
func Open(path string) (*Store, error) {
	_ = schemaSQL // TODO: exec schema on open
	// TODO: sql.Open("sqlite3", path), PRAGMA journal_mode=WAL, PRAGMA foreign_keys=ON
	return &Store{}, nil
}

// Close releases the underlying handle.
func (s *Store) Close() error {
	return nil
}
