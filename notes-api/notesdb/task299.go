package notesdb

import (
	"context"
	"database/sql"
	"fmt"
)

func ListNotes(ctx context.Context, db *sql.DB) ([]Note, error) {
	const query = `SELECT id, title, body, created_at FROM notes ORDER BY id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list notes query: %w", err)
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan note row: %w", err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notes, nil
}
