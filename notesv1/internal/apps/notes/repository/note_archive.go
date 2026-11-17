package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/entities"
	"notesv1/internal/apps/notes/usecases"
)

func (r repo) ArchiveNote(ctx context.Context, q dto.ArchiveNoteQuery) (dto.ArchiveNoteResult, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET archived_at = $2 WHERE id=$1",
		entities.NotesTable,
	)

	result, err := r.db.ExecContext(
		ctx,
		query,
		q.ID,
		q.ArchivedAt,
	)
	if err != nil {
		return dto.ArchiveNoteResult{}, fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return dto.ArchiveNoteResult{}, usecases.ErrNoteNotFound
	}

	return dto.ArchiveNoteResult{Note: nil}, nil
}
