package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r repo) addNoteTags(ctx context.Context, noteID int64, userID uuid.UUID, names []string) error {
	for _, name := range names {
		var tagID uuid.UUID

		err := r.db.QueryRowContext(ctx, `
            INSERT INTO tags (user_id, name)
            VALUES ($1, $2)
            ON CONFLICT (user_id, name)
            DO UPDATE SET name = EXCLUDED.name
            RETURNING id
        `, userID, name).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("get or create tag %q: %w", name, err)
		}

		_, err = r.db.ExecContext(ctx, `
            INSERT INTO note_tags (note_id, tag_id)
            VALUES ($1, $2)
        `, noteID, tagID)
		if err != nil {
			return fmt.Errorf("link tag %q to note: %w", name, err)
		}
	}

	return nil
}

func (r repo) loadNoteTags(ctx context.Context, noteID int64, userID uuid.UUID) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT t.name
        FROM note_tags AS nt
        JOIN tags AS t ON t.id = nt.tag_id
        JOIN notes AS n ON n.id = nt.note_id
        WHERE n.id = $1
          AND n.user_id = $2
          AND t.user_id = $2
        ORDER BY t.name
    `, noteID, userID)
	if err != nil {
		return nil, fmt.Errorf("load note tags: %w", err)
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan tag name: %w", err)
		}
		tags = append(tags, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read tag names: %w", err)
	}

	return tags, nil
}
