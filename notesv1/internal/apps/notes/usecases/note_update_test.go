package usecases_test

import (
	"errors"
	"notesv1/internal/apps/notes/dto"
	"notesv1/internal/apps/notes/usecases"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpdateNote(t *testing.T) {
	t.Parallel()

	t.Run("tc1: ok, updates title and body", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()

		inputID := int64(42)
		inputTitle := "title1"
		inputBody := new("body1")
		inputUserID := uuid.New()
		now := time.Now()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().UpdateNote(ctx, dto.UpdateNoteCommand{
			ID:        inputID,
			Title:     inputTitle,
			Body:      inputBody,
			UserID:    inputUserID,
			Tags:      []string{},
			UpdatedAt: now,
		}).Return(dto.UpdateNoteResult{}, nil).Times(1)

		output, err := newUsecases(deps).UpdateNote(ctx, usecases.UpdateNoteInput{
			ID:     inputID,
			Title:  inputTitle,
			Body:   inputBody,
			UserID: inputUserID,
		})
		r.NoError(err)
		r.Equal(usecases.UpdateNoteOutput{Tags: []string{}}, output)
	})

	t.Run("tc2: ok, nil body", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		now := time.Now()
		userID := uuid.New()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().UpdateNote(ctx, dto.UpdateNoteCommand{
			ID:        42,
			Title:     "title1",
			Body:      nil,
			UserID:    userID,
			Tags:      []string{},
			UpdatedAt: now,
		}).Return(dto.UpdateNoteResult{}, nil)

		output, err := newUsecases(deps).UpdateNote(ctx, usecases.UpdateNoteInput{ID: 42, Title: "title1", UserID: userID})
		r.NoError(err)
		r.Equal(usecases.UpdateNoteOutput{Tags: []string{}}, output)
	})

	t.Run("tc3: ok, title has 256 symbols", func(t *testing.T) {
		t.Parallel()

		r := require.New(t)
		deps := newDeps(t)
		ctx := newCtx()
		title := strings.Repeat("я", 256)
		now := time.Now()
		userID := uuid.New()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().UpdateNote(ctx, dto.UpdateNoteCommand{
			ID:        42,
			Title:     title,
			Body:      nil,
			UserID:    userID,
			Tags:      []string{},
			UpdatedAt: now,
		}).Return(dto.UpdateNoteResult{}, nil)

		output, err := newUsecases(deps).UpdateNote(ctx, usecases.UpdateNoteInput{ID: 42, Title: title, UserID: userID})
		r.NoError(err)
		r.Equal(usecases.UpdateNoteOutput{Tags: []string{}}, output)
	})

	t.Run("tc4: error, zero id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).UpdateNote(newCtx(), usecases.UpdateNoteInput{Title: "title1"})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc5: error, negative id", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).UpdateNote(newCtx(), usecases.UpdateNoteInput{ID: -1, Title: "title1"})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc6: error, empty title", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).UpdateNote(newCtx(), usecases.UpdateNoteInput{ID: 42})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc7: error, title is longer than 256 symbols", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		title := strings.Repeat("я", 257)
		output, err := newUsecases(deps).UpdateNote(newCtx(), usecases.UpdateNoteInput{ID: 42, Title: title, UserID: uuid.New()})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc8: error, empty body", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		output, err := newUsecases(deps).UpdateNote(newCtx(), usecases.UpdateNoteInput{
			ID:     42,
			Title:  "title1",
			Body:   new(""),
			UserID: uuid.New(),
		})
		require.Error(t, err)
		require.ErrorIs(t, err, usecases.ErrValidation)
		require.Zero(t, output)
	})

	t.Run("tc9: error, note not found", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()
		now := time.Now()
		userID := uuid.New()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().UpdateNote(ctx, dto.UpdateNoteCommand{ID: 42, Title: "title1", UserID: userID, Tags: []string{}, UpdatedAt: now}).
			Return(dto.UpdateNoteResult{}, usecases.ErrNoteNotFound)

		output, err := newUsecases(deps).UpdateNote(ctx, usecases.UpdateNoteInput{ID: 42, Title: "title1", UserID: userID})
		require.ErrorIs(t, err, usecases.ErrNoteNotFound)
		require.Zero(t, output)
	})

	t.Run("tc10: error, repository failure", func(t *testing.T) {
		t.Parallel()

		deps := newDeps(t)
		ctx := newCtx()
		repositoryErr := errors.New("database unavailable")
		now := time.Now()
		userID := uuid.New()

		deps.timeProvider.EXPECT().Now().Return(now)
		deps.repo.EXPECT().UpdateNote(ctx, dto.UpdateNoteCommand{ID: 42, Title: "title1", UserID: userID, Tags: []string{}, UpdatedAt: now}).
			Return(dto.UpdateNoteResult{}, repositoryErr)

		output, err := newUsecases(deps).UpdateNote(ctx, usecases.UpdateNoteInput{ID: 42, Title: "title1", UserID: userID})
		require.ErrorIs(t, err, repositoryErr)
		require.Zero(t, output)
	})
}
