package usecases

import (
	"context"
	"notesv1/internal/apps/notes/dto"
	"time"
)

type Repository interface {
	ListNotes(ctx context.Context, q dto.ListNotesQuery) (dto.ListNotesResult, error)
	CreateNote(ctx context.Context, note dto.CreateNoteCommand) (dto.CreateNoteResult, error)
	GetNote(ctx context.Context, q dto.GetNoteQuery) (dto.GetNoteResult, error)
	UpdateNote(ctx context.Context, q dto.UpdateNoteCommand) (dto.UpdateNoteResult, error)
	DeleteNote(ctx context.Context, q dto.DeleteNoteCommand) (dto.DeleteNoteResult, error)
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
	GetRepository() Repository
}

type TimeProvider interface {
	Now() time.Time
}

type Usecases struct {
	transactionManager TransactionManager
	timeProvider       TimeProvider
}

type NewArgs struct {
	TransactionManager TransactionManager
	TimeProvider       TimeProvider
}

func New(args *NewArgs) *Usecases {
	return &Usecases{
		transactionManager: args.TransactionManager,
		timeProvider:       args.TimeProvider,
	}
}
