package repository

import (
	"context"
	"database/sql"
	"fmt"

	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
)

func (r repo) UpdateNote(ctx context.Context, q dto.UpdateNoteCommand) (dto.UpdateNoteResult, error) {
	body := sql.NullString{}
	if q.Body != nil {
		body = sql.NullString{String: *q.Body, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, `
        UPDATE notes
        SET title = $1, body = $2, updated_at = $3
        WHERE id = $4 AND user_id = $5
    `, q.Title, body, q.UpdatedAt, q.ID, q.UserID)
	if err != nil {
		return dto.UpdateNoteResult{}, fmt.Errorf("update note: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return dto.UpdateNoteResult{}, fmt.Errorf("check updated note: %w", err)
	}
	if affected == 0 {
		return dto.UpdateNoteResult{}, usecases.ErrNoteNotFound
	}

	_, err = r.db.ExecContext(ctx, `
        DELETE FROM note_tags WHERE note_id = $1
    `, q.ID)
	if err != nil {
		return dto.UpdateNoteResult{}, fmt.Errorf("remove old note tags: %w", err)
	}

	if err := r.addNoteTags(ctx, q.ID, q.UserID, q.Tags); err != nil {
		return dto.UpdateNoteResult{}, err
	}

	return dto.UpdateNoteResult{}, nil
}
