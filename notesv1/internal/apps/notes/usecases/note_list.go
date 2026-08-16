package usecases

import (
	"context"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
)

type ListNotesInput struct {
}

type ListNotesOutput struct {
	Items []*models.Note
}

func (u *Usecases) ListNotes(ctx context.Context, input ListNotesInput) (ListNotesOutput, error) {
	repo := u.transactionManager.GetRepository()

	listNotesResult, err := repo.ListNotes(ctx, dto.ListNotesQuery{})
	if err != nil {
		return ListNotesOutput{}, fmt.Errorf("failed to list notes: %w", err)
	}

	return ListNotesOutput{
		Items: listNotesResult.Items,
	}, nil
}
