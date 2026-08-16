package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetNoteInput struct {
	ID int64
}

func (in *GetNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.ID, validation.Required, validation.Min(1)),
	)
}

type GetNoteOutput struct {
	Note *models.Note
}

func (u *Usecases) GetNote(ctx context.Context, input GetNoteInput) (GetNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return GetNoteOutput{}, fmt.Errorf("invalid get note input: %w", errors.Join(err, ErrValidation))
	}

	repo := u.transactionManager.GetRepository()

	result, err := repo.GetNote(ctx, dto.GetNoteQuery{ID: input.ID})
	if err != nil {
		return GetNoteOutput{}, fmt.Errorf("failed to get note: %w", err)
	}

	return GetNoteOutput{Note: result.Note}, nil
}
