package usecases_test

import (
	"errors"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/models"
	"notesv1/internal/apps/notes/usecases"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetNote(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, returns note by id", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		inputID := int64(42)
		note := &models.Note{
			ID:        inputID,
			Title:     "title1",
			Body:      new(string),
			CreatedAt: time.Now(),
		}

		deps.repo.EXPECT().GetNote(ctx, dto.GetNoteQuery{
			ID: inputID,
		}).Return(dto.GetNoteResult{Note: note}, nil).Times(1)

		output, err := newUsecases(deps).GetNote(ctx, usecases.GetNoteInput{ID: inputID})
		r.NoError(err)
		r.Equal(note, output.Note)
	})

	t.Run("tc2: error, zero id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).GetNote(newCtx(), usecases.GetNoteInput{})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc3: error, negative id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).GetNote(newCtx(), usecases.GetNoteInput{ID: -1})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc4: error, note not found", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().GetNote(ctx, dto.GetNoteQuery{ID: 42}).
			Return(dto.GetNoteResult{}, usecases.ErrNoteNotFound)

		output, err := newUsecases(deps).GetNote(ctx, usecases.GetNoteInput{ID: 42})
		require.ErrorIs(t, err, usecases.ErrNoteNotFound)
		require.Zero(t, output)
	})

	t.Run("tc5: error, repository failure", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()
		repositoryErr := errors.New("database unavailable")

		deps.repo.EXPECT().GetNote(ctx, dto.GetNoteQuery{ID: 42}).
			Return(dto.GetNoteResult{}, repositoryErr)

		output, err := newUsecases(deps).GetNote(ctx, usecases.GetNoteInput{ID: 42})
		require.ErrorIs(t, err, repositoryErr)
		require.Zero(t, output)
	})
}
