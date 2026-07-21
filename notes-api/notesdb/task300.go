package notesdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Note struct {
	ID          int64
	Title, Body string
	CreatedAt   time.Time
}

var ErrNoteNotFound = errors.New("note not found")

func GetNoteByID(ctx context.Context, db *sql.DB, id int64) (Note, error) {
	var note Note
	const query = `SELECT id, title,body,created_at FROM notes WHERE id=?`
	err := db.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, fmt.Errorf(" row not found %w", ErrNoteNotFound)
		}
		return Note{}, fmt.Errorf("another error: %w", err)
	}
	return note, nil
}
