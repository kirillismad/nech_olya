package usecases

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"

	"github.com/google/uuid"
)

type NoteCountInput struct {
	UserID uuid.UUID
}

type NoteCountOutput struct {
	Count int64
}

func (u *Usecases) NoteCount(ctx context.Context, input NoteCountInput) (NoteCountOutput, error) {
	repo := u.transactionManager.GetRepository()

	result, err := repo.NoteCount(ctx, dto.CountNotesQuery{UserID: input.UserID})
	if err != nil {
		return NoteCountOutput{}, fmt.Errorf("get note count: %w", err)
	}

	return NoteCountOutput{Count: result.Count}, nil
}
