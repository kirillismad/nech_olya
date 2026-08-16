package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DeleteNoteInput struct {
	ID int64
}

func (in *DeleteNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.ID, validation.Required, validation.Min(1)),
	)
}

type DeleteNoteOutput struct{}

func (u *Usecases) DeleteNote(ctx context.Context, input DeleteNoteInput) (DeleteNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return DeleteNoteOutput{}, fmt.Errorf("invalid delete note input: %w", errors.Join(err, ErrValidation))
	}

	repo := u.transactionManager.GetRepository()

	_, err := repo.DeleteNote(ctx, dto.DeleteNoteCommand{ID: input.ID})
	if err != nil {
		return DeleteNoteOutput{}, fmt.Errorf("failed to delete note: %w", err)
	}

	return DeleteNoteOutput{}, nil
}
