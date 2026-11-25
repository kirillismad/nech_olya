package usecases_test

import (
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestoreNote(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, restore note", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().RestoreNote(ctx, dto.RestoreNoteQuery{
			ID: 42,
		}).Return(dto.RestoreNoteResult{}, nil)

		output, err := newUsecases(deps).RestoreNote(ctx, usecases.RestoreNoteInput{
			ID: 42,
		})
		r.NoError(err)
		r.Equal(usecases.RestoreNoteOutput{}, output)
	})

	t.Run("tc2: error, invalid id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)

		output, err := newUsecases(deps).RestoreNote(newCtx(), usecases.RestoreNoteInput{
			ID: 0,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc3: error, note not found", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().RestoreNote(ctx, dto.RestoreNoteQuery{
			ID: 42,
		}).Return(dto.RestoreNoteResult{}, usecases.ErrNoteNotFound)

		output, err := newUsecases(deps).RestoreNote(ctx, usecases.RestoreNoteInput{
			ID: 42,
		})
		require.ErrorIs(t, err, usecases.ErrNoteNotFound)
		require.Zero(t, output)
	})

	t.Run("tc4: ok, ideidempotence", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().RestoreNote(ctx, dto.RestoreNoteQuery{
			ID: 42,
		}).Return(dto.RestoreNoteResult{}, nil)

		output, err := newUsecases(deps).RestoreNote(ctx, usecases.RestoreNoteInput{
			ID: 42,
		})
		r.NoError(err)
		r.Equal(usecases.RestoreNoteOutput{}, output)
	})
}
