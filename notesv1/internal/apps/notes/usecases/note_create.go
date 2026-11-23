package usecases

import (
	"context"
	"errors"
	"fmt"
	"notesv1/internal/apps/notes/dto"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CreateNoteInput struct {
	Title  string
	Body   *string
	UserID uuid.UUID
	Tags   []string
}

func (in *CreateNoteInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(&in.Title, validation.Required, validation.RuneLength(1, 256)),
		validation.Field(&in.Body, validation.NilOrNotEmpty),
	)
}

type CreateNoteOutput struct {
	ID   int64
	Tags []string
}

func (u *Usecases) CreateNote(ctx context.Context, input CreateNoteInput) (CreateNoteOutput, error) {
	if err := input.Validate(); err != nil {
		return CreateNoteOutput{}, fmt.Errorf("invalid create note input: %w", errors.Join(err, ErrValidation))
	}

	if input.UserID == uuid.Nil {
		return CreateNoteOutput{}, ErrValidation
	}

	tags, err := normalizeTags(input.Tags)
	if err != nil {
		return CreateNoteOutput{}, err
	}

	now := u.timeProvider.Now()
	var createdID int64

	err = u.transactionManager.WithTransaction(ctx, func(repo Repository) error {
		result, err := repo.CreateNote(ctx, dto.CreateNoteCommand{
			UserID:    input.UserID,
			Title:     input.Title,
			Body:      input.Body,
			Tags:      tags,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		createdID = result.ID
		return nil
	})
	if err != nil {
		return CreateNoteOutput{}, err
	}

	return CreateNoteOutput{ID: createdID, Tags: tags}, nil
}
