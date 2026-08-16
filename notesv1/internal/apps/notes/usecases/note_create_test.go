package usecases_test

import (
	"errors"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	t.Parallel()
	t.Run("tc1: ok, set title and body", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)

		ctx := newCtx()
		now := time.Now()

		inputTitle := "title1"
		inputBody := new("body1")
		createdID := int64(42)

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().CreateNote(ctx, dto.CreateNoteCommand{
			Title:     inputTitle,
			Body:      inputBody,
			CreatedAt: now,
			UpdatedAt: now,
		}).Return(dto.CreateNoteResult{ID: createdID}, nil).Times(1)

		input := usecases.CreateNoteInput{
			Title: inputTitle,
			Body:  inputBody,
		}
		output, err := newUsecases(deps).CreateNote(ctx, input)
		r.NoError(err)
		r.Equal(createdID, output.ID)
	})

	t.Run("tc2: ok, nil body", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		now := time.Now()
		createdID := int64(43)

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().CreateNote(ctx, dto.CreateNoteCommand{
			Title:     "title1",
			Body:      nil,
			CreatedAt: now,
			UpdatedAt: now,
		}).Return(dto.CreateNoteResult{ID: createdID}, nil)

		output, err := newUsecases(deps).CreateNote(ctx, usecases.CreateNoteInput{Title: "title1"})
		r.NoError(err)
		r.Equal(createdID, output.ID)
	})

	t.Run("tc3: ok, title has 256 symbols", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		now := time.Now()
		createdID := int64(44)
		title := strings.Repeat("я", 256)

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().CreateNote(ctx, dto.CreateNoteCommand{
			Title:     title,
			Body:      nil,
			CreatedAt: now,
			UpdatedAt: now,
		}).Return(dto.CreateNoteResult{ID: createdID}, nil)

		output, err := newUsecases(deps).CreateNote(ctx, usecases.CreateNoteInput{Title: title})
		r.NoError(err)
		r.Equal(createdID, output.ID)
	})

	t.Run("tc4: error, empty title", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).CreateNote(newCtx(), usecases.CreateNoteInput{})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output.ID)
	})

	t.Run("tc5: error, title is longer than 256 symbols", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		title := strings.Repeat("я", 257)
		output, err := newUsecases(deps).CreateNote(newCtx(), usecases.CreateNoteInput{Title: title})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output.ID)
	})

	t.Run("tc6: error, empty body", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).CreateNote(newCtx(), usecases.CreateNoteInput{
			Title: "title1",
			Body:  new(""),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output.ID)
	})

	t.Run("tc7: error, repository failure", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()
		repositoryErr := errors.New("database unavailable")

		deps.timeProvider.EXPECT().Now().Return(time.Now())
		deps.repo.EXPECT().CreateNote(mock.Anything, mock.Anything).Return(dto.CreateNoteResult{}, repositoryErr)

		_, err := newUsecases(deps).CreateNote(ctx, usecases.CreateNoteInput{Title: "title1"})
		require.ErrorIs(t, err, repositoryErr)
	})
}
