package notesdb

import (
	"context"
	"database/sql"
	"errors"
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
	err := db.QueryRowContext(ctx, `
SELECT id, title,body,created_at
FROM notes
WHERE id=?`, id).Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, ErrNoteNotFound
		}
		return Note{}, err
	}
	return note, nil
}
