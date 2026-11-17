package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RestoreNoteInput struct {
	ID int64
}

func (in *RestoreNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.ID, validation.Required, validation.Min(1)),
	)
}

type RestoreNoteOutput struct {
	Note *models.Note
}

func (u *Usecases) RestoreNote(ctx context.Context, input RestoreNoteInput) (RestoreNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return RestoreNoteOutput{}, fmt.Errorf("invalid get note input: %w", errors.Join(err, ErrValidation))
	}

	repo := u.transactionManager.GetRepository()

	result, err := repo.RestoreNote(ctx, dto.RestoreNoteQuery{ID: input.ID})
	if err != nil {
		return RestoreNoteOutput{}, fmt.Errorf("failed to get note: %w", err)
	}

	return RestoreNoteOutput{Note: result.Note}, nil
}
