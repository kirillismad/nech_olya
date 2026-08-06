package notesdb

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
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

func insertNote(t *testing.T, db *sql.DB, title, body string) {
	_, err := db.Exec(`INSERT INTO notes (title, body) VALUES (?, ?)`, title, body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListNotes_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	notes, err := ListNotes(context.Background(), db)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if notes == nil {
		t.Error("slice is nil")
	}
	if len(notes) != 0 {
		t.Errorf("len = %d, want 0", len(notes))
	}
}

func TestListNotes_ReturnsAllInOrder(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	insertNote(t, db, "Title 1", "Body 1")
	insertNote(t, db, "Title 2", "Body 2")
	insertNote(t, db, "Title 3", "Body 3")

	notes, err := ListNotes(context.Background(), db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(notes) != 3 {
		t.Fatalf("len = %d, want 3", len(notes))
	}

	for i := 0; i < len(notes); i++ {
		if notes[i].ID != int64(i+1) {
			t.Errorf("notes[%d].ID = %d, want %d", i, notes[i].ID, i+1)
		}
	}
}

func TestListNotes_FieldsMatch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	insertNote(t, db, "Title A", "Body A")
	insertNote(t, db, "Title B", "Body B")

	notes, err := ListNotes(context.Background(), db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(notes) != 2 {
		t.Fatalf("len = %d, want 2", len(notes))
	}

	if notes[0].Title != "Title A" || notes[0].Body != "Body A" {
		t.Errorf("first note = %+v, want {Title: Title A, Body: Body A}", notes[0])
	}
	if notes[1].Title != "Title B" || notes[1].Body != "Body B" {
		t.Errorf("second note = %+v, want {Title: Title B, Body: Body B}", notes[1])
	}
}
