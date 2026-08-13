package notesrepositorygo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (Note, error) {
	var entity NoteEntity

	const getByIdQuery = `SELECT id,title,body,created_at FROM notes WHERE id=?`
	err := r.db.QueryRowContext(ctx, getByIdQuery, id).Scan(&entity.ID, &entity.Title, &entity.Body, &entity.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, ErrNotFound
		}
		return Note{}, fmt.Errorf("failed get data DB %w", err)
	}
	return entityToNote(entity), nil
}
