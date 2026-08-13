package notesrepositorygo

import (
	"context"
	"fmt"
)

func (r *Repository) List(ctx context.Context) ([]Note, error) {
	const listQuery = `SELECT id,title,body,created_at FROM notes ORDER BY id`
	rows, err := r.db.QueryContext(ctx, listQuery)
	if err != nil {
		return nil, fmt.Errorf("failed get notes %w", err)
	}
	defer rows.Close()
	notes := make([]Note, 0)
	for rows.Next() {
		var entity NoteEntity
		err := rows.Scan(&entity.ID, &entity.Title, &entity.Body, &entity.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert the data %w", err)
		}
		notes = append(notes, entityToNote(entity))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error data transformations %w", err)
	}
	return notes, nil
}
