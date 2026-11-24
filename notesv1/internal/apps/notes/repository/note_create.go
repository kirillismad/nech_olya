package repository

import (
	"context"
	"database/sql"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/entities"
)

func (r repo) CreateNote(ctx context.Context, q dto.CreateNoteCommand) (dto.CreateNoteResult, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5) RETURNING %s",
		entities.NotesTable,
		entities.NotesTitleColumn,
		entities.NotesBodyColumn,
		entities.NotesCreatedAtColumn,
		entities.NotesUpdatedAtColumn,
		entities.NotesUserIDColumn,
		entities.NotesIDColumn,
	)
	entity := entities.Note{
		Title: q.Title,
		Body: func() sql.NullString {
			if q.Body == nil {
				return sql.NullString{}
			}

			return sql.NullString{String: *q.Body, Valid: true}
		}(),
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
		UserID:    q.UserID,
	}

	var id int64

	err := r.db.QueryRowContext(
		ctx,
		query,
		entity.Title,
		entity.Body,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.UserID,
	).Scan(&id)
	if err != nil {
		return dto.CreateNoteResult{}, fmt.Errorf("failed to create note: %w", err)
	}

	return dto.CreateNoteResult{ID: id}, nil
}
