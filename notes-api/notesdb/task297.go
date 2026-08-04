package notesdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

//Объявить структуру NullableNote { ID int64; Title string; Body sql.NullString; CreatedAt time.Time } и функции:
// - InsertNullable(ctx, db, title string, body *string) (int64, error) — если body == nil, вставляет NULL, иначе — значение
// - GetNullable(ctx, db, id int64) (NullableNote, error) — корректно сканирует body через sql.NullString

// Тесты (in-memory + Migrate):
// - TestInsertNullable_WithBody — после InsertNullable(ctx, db, "t", strPtr("hello")) GetNullable возвращает Body.Valid == true и Body.String == "hello"
// - TestInsertNullable_NilBody — после InsertNullable(ctx, db, "t", nil) GetNullable возвращает Body.Valid == false
// - TestScanWithoutNullableFails — попытка db.QueryRow("SELECT body FROM notes WHERE id = ?", idWithNull).Scan(&s) в обычный string возвращает ошибку (демонстрирует, зачем нужен sql.NullString)

type NullableNote struct {
	ID        int64
	Title     string
	Body      sql.NullString
	CreatedAt time.Time
}

func InsertNullable(ctx context.Context, db *sql.DB, title string, body *string) (int64, error) {
	const query = `INSERT INTO notes (title, body) VALUES (?, ?)`

	nullableNote := NullableNote{
		Title: title,
	}
	if body != nil {
		nullableNote.Body = sql.NullString{String: *body, Valid: true}
	}

	result, err := db.ExecContext(ctx, query, nullableNote.Title, nullableNote.Body)
	if err != nil {
		return 0, fmt.Errorf("changes to the notes %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to extract id %w", err)
	}
	return id, nil
}

func GetNullable(ctx context.Context, db *sql.DB, id int64) (NullableNote, error) {
	const query = `SELECT id, title, body, created_at FROM notes WHERE id = ?`

	var note NullableNote
	err := db.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return NullableNote{}, fmt.Errorf("Failed to extract the string %w", err)
	}
	return note, nil
}
