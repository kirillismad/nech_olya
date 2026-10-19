package usecases

import (
	"context"
	"fmt"
)

type NoteCountOutput struct {
	Count int64
}

func (u *Usecases) NoteCount(ctx context.Context) (NoteCountOutput, error) {
	repo := u.transactionManager.GetRepository()

	result, err := repo.NoteCount(ctx)
	if err != nil {
		return NoteCountOutput{}, fmt.Errorf("get note count: %w", err)
	}
	return NoteCountOutput{Count: result.Count}, nil
}
