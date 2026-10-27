package usecases

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
)

type ListNotesInput struct {
	Limit  int64
	Offset int64
}

func(input ListNotesInput) Validate() error {
if input.Limit < 1 || input.Limit > 100 {
		return ErrValidation
	}

	if input.Offset < 0 {
		return ErrValidation
	}

	return nil
}

type ListNotesOutput struct {
	Items []*models.Note
}

func (u *Usecases) ListNotes(ctx context.Context, input ListNotesInput) (ListNotesOutput, error) {
	if err := input.Validate(); err != nil {
		return ListNotesOutput{}, err
	}
	repo := u.transactionManager.GetRepository()

	listNotesResult, err := repo.ListNotes(ctx, dto.ListNotesQuery{Limit: input.Limit, Offset: input.Offset})
	if err != nil {
		return ListNotesOutput{}, fmt.Errorf("failed to list notes: %w", err)
	}

	return ListNotesOutput{
		Items: listNotesResult.Items,
	}, nil
}
