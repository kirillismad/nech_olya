package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
	"notesv1/internal/apps/notes/usecases"
	"notesv1/internal/entities"
)

func (r repo) GetNote(ctx context.Context, q dto.GetNoteQuery) (dto.GetNoteResult, error) {
	query := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s FROM %s WHERE %s = ?`,
		entities.NotesIDColumn,
		entities.NotesTitleColumn,
		entities.NotesBodyColumn,
		entities.NotesCreatedAtColumn,
		entities.NotesUpdatedAtColumn,
		entities.NotesTable,
		entities.NotesIDColumn,
	)

	var note entities.Note
	err := r.db.QueryRowContext(ctx, query, q.ID).Scan(
		&note.ID,
		&note.Title,
		&note.Body,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GetNoteResult{}, errors.Join(usecases.ErrNoteNotFound, sql.ErrNoRows)
		}
		return dto.GetNoteResult{}, fmt.Errorf("failed to get note: %w", err)
	}

	return dto.GetNoteResult{Note: &models.Note{
		ID:    note.ID,
		Title: note.Title,
		Body: func() *string {
			var body sql.NullString = note.Body
			if !body.Valid {
				return nil
			}
			return &body.String
		}(),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}}, nil
}
