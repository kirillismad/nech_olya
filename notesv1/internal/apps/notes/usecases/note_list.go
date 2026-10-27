package usecases

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
)

type ListNotesInput struct {
	Limit  int
	Offset int
}

type ListNotesOutput struct {
	Items []*models.Note
}

func (u *Usecases) ListNotes(ctx context.Context, input ListNotesInput) (ListNotesOutput, error) {
	repo := u.transactionManager.GetRepository()

	if input.Limit < 0 || input.Limit > 100 {
		return ListNotesOutput{}, ErrInvalidLimit
	}

	if input.Offset < 0 {
		return ListNotesOutput{}, ErrInvalidOffset
	}

	listNotesResult, err := repo.ListNotes(ctx, dto.ListNotesQuery{Limit: input.Limit, Offset: input.Offset})
	if err != nil {
		return ListNotesOutput{}, fmt.Errorf("failed to list notes: %w", err)
	}

	return ListNotesOutput{
		Items: listNotesResult.Items,
	}, nil
}
