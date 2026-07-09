package notesdb

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func OpenInMemory() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to Ping: %w", err)
	}

	return db, nil
}

