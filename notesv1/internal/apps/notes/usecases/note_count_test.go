package usecases_test

import (
	"errors"
	"notesv1/internal/apps/notes/dto"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoteCount(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, returns note count", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		expectedCount := int64(42)
		deps.repo.EXPECT().NoteCount(ctx).Return(dto.CountNotesResult{Count: expectedCount}, nil)

		output, err := newUsecases(deps).NoteCount(ctx)
		r.NoError(err)
		r.Equal(expectedCount, output)
	})

	t.Run("tc2: error, returns error", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		repositoryErr := errors.New("repository error")
		deps.repo.EXPECT().NoteCount(ctx).Return(dto.CountNotesResult{}, repositoryErr)

		output, err := newUsecases(deps).NoteCount(ctx)
		r.Error(err)
		r.Zero(output)
	})
}
