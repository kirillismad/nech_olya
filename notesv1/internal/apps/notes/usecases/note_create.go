package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateNoteInput struct {
	Title string
	Body  *string
}

func (in *CreateNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Title, validation.Required, validation.RuneLength(1, 256)),
		validation.Field(&in.Body, validation.NilOrNotEmpty),
	)
}

type CreateNoteOutput struct {
	ID int64
}

func (u *Usecases) CreateNote(ctx context.Context, input CreateNoteInput) (CreateNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return CreateNoteOutput{}, fmt.Errorf("invalid create note input: %w", errors.Join(err, ErrValidation))
	}

	now := u.timeProvider.Now()

	repo := u.transactionManager.GetRepository()

	result, err := repo.CreateNote(ctx, dto.CreateNoteCommand{
		Title:     input.Title,
		Body:      input.Body,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return CreateNoteOutput{}, fmt.Errorf("failed to create note: %w", err)
	}

	return CreateNoteOutput{ID: result.ID}, nil
}
