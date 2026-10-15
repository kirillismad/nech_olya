package usecases

import (
	"context"
)

func (u *Usecases) NoteCount(ctx context.Context) (int64, error) {
	repo := u.transactionManager.GetRepository()

	result, err := repo.NoteCount(ctx)
	if err != nil {
		return 0, err
	}

	return result.Count, nil
}
