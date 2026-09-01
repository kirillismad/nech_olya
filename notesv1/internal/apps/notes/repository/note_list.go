package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
	"notesv1/internal/entities"
)

func (r repo) ListNotes(ctx context.Context, q dto.ListNotesQuery) (dto.ListNotesResult, error) {
	query := fmt.Sprintf(
		"SELECT %s, %s, %s, %s, %s FROM %s",
		entities.NotesIDColumn,
		entities.NotesTitleColumn,
		entities.NotesBodyColumn,
		entities.NotesCreatedAtColumn,
		entities.NotesUpdatedAtColumn,
		entities.NotesTable,
	)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return dto.ListNotesResult{}, fmt.Errorf("failed to list notes: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var entityList []entities.Note

	for rows.Next() {
		var item entities.Note
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Body,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return dto.ListNotesResult{}, fmt.Errorf("failed to scan note: %w", err)
		}

		entityList = append(entityList, item)
	}

	if err := rows.Err(); err != nil {
		return dto.ListNotesResult{}, fmt.Errorf("rows error: %w", err)
	}

	modelList := make([]*models.Note, 0, len(entityList))
	for _, item := range entityList {
		model := &models.Note{
			ID:    item.ID,
			Title: item.Title,
			Body: func() *string {
				if !item.Body.Valid {
					return nil
				}

				return &item.Body.String
			}(),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		modelList = append(modelList, model)
	}

	return dto.ListNotesResult{
		Items: modelList,
	}, nil
}
