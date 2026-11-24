package repository

import (
	"context"
	"database/sql"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/entities"
)

func (r repo) UpdateNote(ctx context.Context, q dto.UpdateNoteCommand) (dto.UpdateNoteResult, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET %s = $1, %s = $2, %s = $3 WHERE %s = $4 AND %s=$5",
		entities.NotesTable,
		entities.NotesTitleColumn,
		entities.NotesBodyColumn,
		entities.NotesUpdatedAtColumn,
		entities.NotesIDColumn,
		entities.NotesUserIDColumn,
	)

	entity := entities.Note{
		ID:    q.ID,
		Title: q.Title,
		Body: func() sql.NullString {
			if q.Body == nil {
				return sql.NullString{}
			}

			return sql.NullString{String: *q.Body, Valid: true}
		}(),
		UpdatedAt: q.UpdatedAt,
		UserID:    q.UserID,
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		entity.Title,
		entity.Body,
		entity.UpdatedAt,
		entity.ID,
		entity.UserID,
	)
	if err != nil {
		return dto.UpdateNoteResult{}, fmt.Errorf("failed to update note: %w", err)
	}

	return dto.UpdateNoteResult{}, nil
}
