package notesdb

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func OpenInMemory() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
