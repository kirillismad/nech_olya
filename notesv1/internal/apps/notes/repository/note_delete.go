package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/entities"
)

func (r repo) DeleteNote(ctx context.Context, q dto.DeleteNoteCommand) (dto.DeleteNoteResult, error) {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = $1 AND %s = $2",
		entities.NotesTable,
		entities.NotesIDColumn,
		entities.NotesUserIDColumn,
	)

	_, err := r.db.ExecContext(ctx, query, q.ID, q.UserID)
	if err != nil {
		return dto.DeleteNoteResult{}, fmt.Errorf("failed to delete note: %w", err)
	}

	return dto.DeleteNoteResult{}, nil
}
