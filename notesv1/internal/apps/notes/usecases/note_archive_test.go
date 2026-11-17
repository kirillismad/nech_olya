package usecases_test

import (
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestArchiveNote(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, archive note", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		now := time.Now()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().ArchiveNote(ctx, dto.ArchiveNoteQuery{
			ID:         42,
			ArchivedAt: &now,
		}).Return(dto.ArchiveNoteResult{}, nil)

		output, err := newUsecases(deps).ArchiveNote(ctx, usecases.ArchiveNoteInput{
			ID: 42,
		})
		r.NoError(err)
		r.Equal(usecases.ArchiveNoteOutput{}, output)
	})

	t.Run("tc2: error, invalid id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)

		output, err := newUsecases(deps).ArchiveNote(newCtx(), usecases.ArchiveNoteInput{
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
		now := time.Now()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().ArchiveNote(ctx, dto.ArchiveNoteQuery{
			ID:         42,
			ArchivedAt: &now,
		}).Return(dto.ArchiveNoteResult{}, usecases.ErrNoteNotFound)

		output, err := newUsecases(deps).ArchiveNote(ctx, usecases.ArchiveNoteInput{
			ID: 42,
		})
		require.ErrorIs(t, err, usecases.ErrNoteNotFound)
		require.Zero(t, output)
	})
}
