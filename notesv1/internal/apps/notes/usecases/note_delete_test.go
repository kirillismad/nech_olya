package usecases_test

import (
	"errors"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteNote(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, deletes note by id", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		inputID := int64(42)

		deps.repo.EXPECT().DeleteNote(ctx, dto.DeleteNoteCommand{
			ID: inputID,
		}).Return(dto.DeleteNoteResult{}, nil).Times(1)

		output, err := newUsecases(deps).DeleteNote(ctx, usecases.DeleteNoteInput{ID: inputID})
		r.NoError(err)
		r.Equal(usecases.DeleteNoteOutput{}, output)
	})

	t.Run("tc2: error, zero id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).DeleteNote(newCtx(), usecases.DeleteNoteInput{})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc3: error, negative id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).DeleteNote(newCtx(), usecases.DeleteNoteInput{ID: -1})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc4: error, note not found", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().DeleteNote(ctx, dto.DeleteNoteCommand{ID: 42}).
			Return(dto.DeleteNoteResult{}, usecases.ErrNoteNotFound)

		output, err := newUsecases(deps).DeleteNote(ctx, usecases.DeleteNoteInput{ID: 42})
		require.ErrorIs(t, err, usecases.ErrNoteNotFound)
		require.Zero(t, output)
	})

	t.Run("tc5: error, repository failure", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()
		repositoryErr := errors.New("database unavailable")

		deps.repo.EXPECT().DeleteNote(ctx, dto.DeleteNoteCommand{ID: 42}).
			Return(dto.DeleteNoteResult{}, repositoryErr)

		output, err := newUsecases(deps).DeleteNote(ctx, usecases.DeleteNoteInput{ID: 42})
		require.ErrorIs(t, err, repositoryErr)
		require.Zero(t, output)
	})
}
