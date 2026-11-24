package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
	"notesv1/internal/entities"
)

func (r repo) ListNotes(ctx context.Context, q dto.ListNotesQuery) (dto.ListNotesResult, error) {
	var pattern *string
	if q.Query != nil {
		p := "%" + *q.Query + "%"
		pattern = &p
	}
	query := `
    SELECT n.id, n.title, n.body, n.created_at, n.updated_at
    FROM notes AS n
    WHERE n.user_id = $1
      AND (
        $2::text IS NULL
        OR EXISTS (
            SELECT 1
            FROM note_tags AS nt
            JOIN tags AS t ON t.id = nt.tag_id
            WHERE nt.note_id = n.id
              AND t.user_id = n.user_id
              AND t.name = $2
        )
      )
	 AND (
        $3::text IS NULL
        OR n.title ILIKE $3
        OR COALESCE(n.body, '') ILIKE $3
      )
    ORDER BY n.id
    LIMIT $3 OFFSET $4
`

	rows, err := r.db.QueryContext(ctx, query, q.UserID, q.Tag, pattern, q.Limit, q.Offset)
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
		tags, err := r.loadNoteTags(ctx, item.ID, q.UserID)
		if err != nil {
			return dto.ListNotesResult{}, fmt.Errorf("failed to load note tags: %w", err)
		}

		model := &models.Note{
			ID:    item.ID,
			Title: item.Title,
			Body: func() *string {
				if !item.Body.Valid {
					return nil
				}

				return &item.Body.String
			}(),
			Tags:      tags,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		modelList = append(modelList, model)
	}

	return dto.ListNotesResult{
		Items: modelList,
	}, nil
}
