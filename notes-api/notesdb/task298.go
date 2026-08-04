package notesdb

import (
	"context"
	"database/sql"
	"fmt"
)

func SearchNotes(ctx context.Context, db *sql.DB, pattern string) ([]Note, error) {
	const query = `SELECT id, title, body, created_at FROM notes WHERE title LIKE ? ORDER BY id`
	rows, err := db.QueryContext(ctx, query, pattern)
	if err != nil {
		return nil, fmt.Errorf("list notes query: %w", err)
	}
	defer rows.Close()
	notes := make([]Note, 0)
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan note row: %w", err)
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return notes, nil
}
