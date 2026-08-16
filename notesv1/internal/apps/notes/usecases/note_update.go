package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateNoteInput struct {
	ID    int64
	Title string
	Body  *string
}

func (in *UpdateNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.ID, validation.Required, validation.Min(1)),
		validation.Field(&in.Title, validation.Required, validation.RuneLength(1, 256)),
		validation.Field(&in.Body, validation.NilOrNotEmpty),
	)
}

type UpdateNoteOutput struct{}

func (u *Usecases) UpdateNote(ctx context.Context, input UpdateNoteInput) (UpdateNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return UpdateNoteOutput{}, fmt.Errorf("invalid update note input: %w", errors.Join(err, ErrValidation))
	}

	now := u.timeProvider.Now()

	err := u.transactionManager.WithTransaction(ctx, func(repo Repository) error {
		note, err := repo.GetNote(ctx, dto.GetNoteQuery{ID: input.ID})
		if err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		_, err = repo.UpdateNote(ctx, dto.UpdateNoteCommand{
			ID:        note.Note.ID,
			Title:     input.Title,
			Body:      input.Body,
			UpdatedAt: now,
		})
		if err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}

		return nil
	})
	if err != nil {
		return UpdateNoteOutput{}, err
	}

	return UpdateNoteOutput{}, nil
}
