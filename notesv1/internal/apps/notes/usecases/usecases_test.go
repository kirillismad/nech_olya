package usecases_test

import (
	"context"
	"notesv1/internal/apps/notes/usecases"
	mock_usecases "notesv1/tests/mocks/apps/notes/usecases"
	"testing"

	"github.com/stretchr/testify/mock"
)

type deps struct {
	transactionManager *mock_usecases.MockTransactionManager
	timeProvider       *mock_usecases.MockTimeProvider
	repo               *mock_usecases.MockRepository
}

func newDeps(t *testing.T) *deps {
	repo := mock_usecases.NewMockRepository(t)
	tm := mock_usecases.NewMockTransactionManager(t)

	tm.EXPECT().GetRepository().Return(repo).Maybe()
	tm.EXPECT().WithTransaction(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(r usecases.Repository) error) error {
		return fn(repo)
	}).Maybe()

	return &deps{
		transactionManager: tm,
		timeProvider:       mock_usecases.NewMockTimeProvider(t),
		repo:               repo,
	}
}

func newUsecases(deps *deps) *usecases.Usecases {
	return usecases.New(&usecases.NewArgs{
		TransactionManager: deps.transactionManager,
		TimeProvider:       deps.timeProvider,
	})
}

func newCtx() context.Context {
	return context.Background()
}
