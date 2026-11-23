package repository

import (
	"context"
	"database/sql"
	"fmt"

	"notesv1/internal/apps/notes/dto"
)

func (r repo) CreateNote(ctx context.Context, q dto.CreateNoteCommand) (dto.CreateNoteResult, error) {
	body := sql.NullString{}
	if q.Body != nil {
		body = sql.NullString{String: *q.Body, Valid: true}
	}

	var noteID int64
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO notes (user_id, title, body, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, q.UserID, q.Title, body, q.CreatedAt, q.UpdatedAt).Scan(&noteID)
	if err != nil {
		return dto.CreateNoteResult{}, fmt.Errorf("create note: %w", err)
	}

	if err := r.addNoteTags(ctx, noteID, q.UserID, q.Tags); err != nil {
		return dto.CreateNoteResult{}, err
	}

	return dto.CreateNoteResult{ID: noteID}, nil
}
