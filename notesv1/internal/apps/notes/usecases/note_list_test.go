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

		deps.repo.EXPECT().ListNotes(ctx, dto.ListNotesQuery{Limit: 10, Offset: 1}).
			Return(dto.ListNotesResult{Items: items}, nil).Times(1)

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{Limit: 10, Offset: 1})
		r.NoError(err)
		r.Equal(items, output.Items)
	})

	t.Run("tc2: ok, returns empty list", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		deps.repo.EXPECT().ListNotes(ctx, dto.ListNotesQuery{Limit: 10, Offset: 1}).
			Return(dto.ListNotesResult{Items: []*models.Note{}}, nil).Times(1)

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{Limit: 10, Offset: 1})
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

	t.Run("tc4: error, invalid limit", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{Limit: -1})
		require.ErrorIs(t, err, usecases.ErrInvalidLimit)
		require.Zero(t, output)
	})

	t.Run("tc5: error, invalid offset", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()

		output, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{Offset: -1})
		require.ErrorIs(t, err, usecases.ErrInvalidOffset)
		require.Zero(t, output)
	})

	t.Run("tc6: ok, limit offset pagination", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()

		expected := dto.ListNotesQuery{
			Limit:  2,
			Offset: 1,
		}

		deps.repo.EXPECT().
			ListNotes(ctx, expected).
			Return(dto.ListNotesResult{
				Items: []*models.Note{
					{ID: 2, Title: "n2"},
					{ID: 3, Title: "n3"},
				},
			}, nil)

		out, err := newUsecases(deps).ListNotes(ctx, usecases.ListNotesInput{
			Limit:  2,
			Offset: 1,
		})

		require.NoError(t, err)
		require.Len(t, out.Items, 2)
		require.Equal(t, int64(2), out.Items[0].ID)
		require.Equal(t, int64(3), out.Items[1].ID)
	})
}
