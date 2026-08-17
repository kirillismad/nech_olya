package notesrepositorygo

import (
	"context"
	"fmt"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const deleteQuery = `DELETE FROM notes WHERE id=?`
	result, err := r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed delete data %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve the number of rows %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
