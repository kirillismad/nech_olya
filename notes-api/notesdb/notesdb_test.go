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
