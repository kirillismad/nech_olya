package notesrepositorygo

import (
	"context"
	"fmt"
)

func (r *Repository) Create(ctx context.Context, note Note) (int64, error) {
	const createQuery = `INSERT INTO notes(title, body) VALUES(?,?)`
	if note.Title == "" {
		return 0, ErrEmptyTitle
	}
	result, err := r.db.ExecContext(ctx, createQuery, note.Title, note.Body)
	if err != nil {
		return 0, fmt.Errorf("failed add a note %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed get %w", err)
	}
	return id, nil
}
