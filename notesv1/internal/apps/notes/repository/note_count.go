package repository

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
)

func (r repo) NoteCount(ctx context.Context, q dto.CountNotesQuery) (dto.CountNotesResult, error) {
	var count int64

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM notes WHERE user_id = $1", q.UserID).Scan(&count)
	if err != nil {
		return dto.CountNotesResult{}, fmt.Errorf("failed to count notes: %w", err)
	}

	return dto.CountNotesResult{Count: count}, nil
}
