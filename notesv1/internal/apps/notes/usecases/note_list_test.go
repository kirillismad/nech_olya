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

func TestListNotes(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, returns notes", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		items := []*models.Note{
			{ID: 1, Title: "title1", Body: new(string), CreatedAt: time.Now()},
			{ID: 2, Title: "title2", CreatedAt: time.Now()},
		}

		deps.repo.EXPECT().ListNotes(ctx, dto.ListNotesQuery{}).
			Return(dto.ListNotesResult{Items: items}, nil).Times(1)

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{})
		r.NoError(err)
		r.Equal(items, output.Items)
	})

	t.Run("tc2: ok, returns empty list", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().ListNotes(ctx, dto.ListNotesQuery{}).
			Return(dto.ListNotesResult{}, nil).Times(1)

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{})
		r.NoError(err)
		r.Empty(output.Items)
	})

	t.Run("tc3: error, repository failure", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()
		repositoryErr := errors.New("database unavailable")

		deps.repo.EXPECT().ListNotes(ctx, dto.ListNotesQuery{}).
			Return(dto.ListNotesResult{}, repositoryErr)

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{})
		require.ErrorIs(t, err, repositoryErr)
		require.Zero(t, output)
	})
}
