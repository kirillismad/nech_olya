package notesrepositorygo

import (
	"context"
	"fmt"
)

func (r *Repository) Update(ctx context.Context, id int64, note Note) error {
	if note.Title == "" {
		return ErrEmptyTitle
	}

	const updateQuery = `UPDATE notes SET title=?,body=? WHERE id=?`
	result, err := r.db.ExecContext(ctx, updateQuery, note.Title, note.Body, id)
	if err != nil {
		return fmt.Errorf("failed update DB %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
