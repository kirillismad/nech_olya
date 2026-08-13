package notesrepositorygo

import (
	"context"
	"fmt"
)

func (r *Repository) BulkCreate(ctx context.Context, notes []Note) ([]int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed start transaction %w", err)
	}
	defer tx.Rollback()
	ids := make([]int64, 0, len(notes))
	const bulkCreateQuery = `INSERT INTO notes(title,body) VALUES(?,?)`
	for _, note := range notes {
		if note.Title == "" {
			return nil, ErrEmptyTitle
		}
		result, err := tx.ExecContext(ctx, bulkCreateQuery, note.Title, note.Body)
		if err != nil {
			return nil, fmt.Errorf("failed create data in DB %w", err)
		}
		affected, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed get id %w", err)
		}
		ids = append(ids, affected)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction rollback %w", err)
	}
	return ids, nil
}
