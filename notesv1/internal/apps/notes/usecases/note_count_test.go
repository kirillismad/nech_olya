package usecases_test

import (
	"errors"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNoteCount(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, returns note count", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		userID := uuid.New()

		expectedCount := int64(42)
		deps.repo.EXPECT().NoteCount(ctx, dto.CountNotesQuery{UserID: userID}).Return(dto.CountNotesResult{Count: expectedCount}, nil)

		output, err := newUsecases(deps).NoteCount(ctx, usecases.NoteCountInput{UserID: userID})
		r.NoError(err)
		r.Equal(usecases.NoteCountOutput{Count: expectedCount}, output)
	})

	t.Run("tc2: error, returns error", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		userID := uuid.New()

		repositoryErr := errors.New("repository error")
		deps.repo.EXPECT().NoteCount(ctx, dto.CountNotesQuery{UserID: userID}).Return(dto.CountNotesResult{}, repositoryErr)

		output, err := newUsecases(deps).NoteCount(ctx, usecases.NoteCountInput{UserID: userID})
		r.Error(err)
		r.Zero(output)
	})
}
