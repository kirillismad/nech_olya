package notesrepositorygo

import (
	"context"
	"fmt"
)

func (r *Repository) Migrate(ctx context.Context) error {
	const migrateQuery = `CREATE TABLE IF NOT EXISTS notes(
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	title TEXT NOT NULL, body TEXT NULL, 
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);`

	_, err := r.db.ExecContext(ctx, migrateQuery)
	if err != nil {
		return fmt.Errorf("failed create table %w", err)
	}
	return nil
}
