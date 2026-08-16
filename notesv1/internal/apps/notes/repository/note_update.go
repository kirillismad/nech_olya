package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/entities"
)

func (r repo) UpdateNote(ctx context.Context, q dto.UpdateNoteCommand) (dto.UpdateNoteResult, error) {
	query := fmt.Sprintf(
		`UPDATE %s SET %s = ?, %s = ?, %s = ? WHERE %s = ?`,
		entities.NotesTable,
		entities.NotesTitleColumn,
		entities.NotesBodyColumn,
		entities.NotesUpdatedAtColumn,
		entities.NotesIDColumn,
	)

	_, err := r.db.ExecContext(ctx, query, q.Title, q.Body, q.UpdatedAt, q.ID)
	if err != nil {
		return dto.UpdateNoteResult{}, fmt.Errorf("failed to update note: %w", err)
	}

	return dto.UpdateNoteResult{}, nil
}
