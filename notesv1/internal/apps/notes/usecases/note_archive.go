package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ArchiveNoteInput struct {
	ID int64
}

func (in *ArchiveNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.ID, validation.Required, validation.Min(1)),
	)
}

type ArchiveNoteOutput struct {
	Note *models.Note
}

func (u *Usecases) ArchiveNote(ctx context.Context, input ArchiveNoteInput) (ArchiveNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return ArchiveNoteOutput{}, fmt.Errorf("invalid get note input: %w", errors.Join(err, ErrValidation))
	}

	now := u.timeProvider.Now()

	repo := u.transactionManager.GetRepository()

	result, err := repo.ArchiveNote(ctx, dto.ArchiveNoteQuery{ID: input.ID, ArchivedAt: &now})
	if err != nil {
		return ArchiveNoteOutput{}, fmt.Errorf("failed to get note: %w", err)
	}

	return ArchiveNoteOutput{Note: result.Note}, nil
}
