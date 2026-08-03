package notesdb

import (
	"context"
	"database/sql"
	"testing"
)

func setupTestDB_WithNull(t *testing.T) *sql.DB {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE notes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            body TEXT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestInsertNullable_WithBody(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	strPtr := "hello"
	id, err := InsertNullable(ctx, "t", &strPtr, db)
	if err != nil {
		t.Fatal(err)
	}
	note, err := GetNullable(ctx, id, db)
	if err != nil {
		t.Fatal(err)
	}
	if note.Body.Valid != true {
		t.Fatal("expected valid value")
	}

	if note.Body.String != "hello" {
		t.Fatal("expected value body hello")
	}
	t.Log(note)

}

func TestInsertNullable_NilBody(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	var shear []int64
	for i := 0; i < 3; i++ {
		id, err := InsertNote(ctx, db, "title", "body")
		if err != nil {
			t.Fatal(err)
		}
		shear = append(shear, id)
	}
	if len(shear) != 3 {
		t.Fatal("incorrect length")
	}
	t.Log(len(shear))
}

func TestScanWithoutNullableFails(t *testing.T) {
	db := setupTestDB_WithNull(t)
	ctx := context.Background()
	id, err := InsertNullable(ctx, "title_null", nil, db)
	if err != nil {
		t.Fatal(err)
	}
	var str string
	err = db.QueryRowContext(ctx, "SELECT body FROM notes WHERE id = ?", id).Scan(&str)
	if err == nil {
		t.Fatal(err)
	}
	t.Log(str)
}
