package notesdb

import (
	"context"
	"database/sql"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestOpenInMemory_PingOk(t *testing.T) {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}

func TestOpenInMemory_IsUsable(t *testing.T) {
	db := newTestDB(t)
	_, err := db.Exec("CREATE TABLE t(id INTEGER)")

	if err != nil {
		t.Fatal(err)
	}
}

func TestMigrate_CreatesTable(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	err := Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	var num int
	err = db.QueryRowContext(ctx, "SELECT count(*) FROM notes").Scan(&num)
	if err != nil {
		t.Fatal(err)
	}
	expected := 0
	if num != expected {
		t.Fatal("expected zero")
	}
}

func TestMigrate_Idempotent(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	err := Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	err = Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMigrate_SchemaColumns(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	err := Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	notes := map[string]bool{
		"id":         false,
		"title":      false,
		"body":       false,
		"created_at": false,
	}

	rows, err := db.QueryContext(ctx, `SELECT name FROM pragma_table_info('notes')`)
	if err != nil {
		t.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		if _, ok := notes[name]; ok {
			notes[name] = true
		}

	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	for _, exist := range notes {
		if !exist {
			t.Fatal(err)
		}
	}
}

func TestInsertNote_ReturnsID(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	err := Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	title := "title"
	body := "body"
	id, err := InsertNote(ctx, db, title, body)
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 {
		t.Fatalf("id: %d", id)
	}
	t.Logf("id: %d ", id)
	id, err = InsertNote(ctx, db, title, body)
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 {
		t.Fatalf("id: %d", id)
	}
	t.Logf("id: %d", id)
}

func TestInsertNote_RowExists(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	err := Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	title := "title"
	body := "body"
	id, err := InsertNote(ctx, db, title, body)
	if err != nil {
		t.Fatal(err)
	}

	var notes []string

	rows, err := db.QueryContext(ctx, `SELECT title, body FROM notes WHERE id = ?`, id)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var title string
		var body string
		if err := rows.Scan(&title, &body); err != nil {
			t.Fatal(err)
		}
		notes = append(notes, title, body)
	}
	for _, note := range notes {
		t.Log(note)
	}
}

func TestInsertNote_SpecialChars(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	if err := Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}

	title := "Robert'); DROP TABLE notes;--"
	body := "test body"

	_, err := InsertNote(ctx, db, title, body)
	if err != nil {
		t.Fatal(err)
	}

	var count int

	err = db.QueryRowContext(
		ctx,
		"SELECT count(*) FROM notes",
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
}
