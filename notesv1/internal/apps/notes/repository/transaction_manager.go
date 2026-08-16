package repository

import (
	"context"
	"database/sql"
	"fmt"
	"notesv1/internal/apps/notes/usecases"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (m *TransactionManager) WithTransaction(ctx context.Context, fn func(repo usecases.Repository) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	r := repo{
		db: tx,
	}
	if err := fn(r); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *TransactionManager) GetRepository() usecases.Repository {
	return repo{
		db: m.db,
	}
}
