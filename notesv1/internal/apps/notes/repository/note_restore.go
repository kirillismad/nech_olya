package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/entities"
	"notesv1/internal/apps/notes/usecases"
)

func (r repo) RestoreNote(ctx context.Context, q dto.RestoreNoteQuery) (dto.RestoreNoteResult, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET archived_at = NULL WHERE id=$1",
		entities.NotesTable,
	)

	result, err := r.db.ExecContext(
		ctx,
		query,
		q.ID,
	)
	if err != nil {
		return dto.RestoreNoteResult{}, fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return dto.RestoreNoteResult{}, usecases.ErrNoteNotFound
	}

	return dto.RestoreNoteResult{Note: nil}, nil
}
