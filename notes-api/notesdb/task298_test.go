package notesdb

import (
	"context"
	"database/sql"
	"testing"
)

func setupTestDb(t *testing.T) *sql.DB {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE notes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            body TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
    INSERT INTO notes (title, body) VALUES 
        ('go basics', 'Learning Go basics'),
        ('go advanced', 'Advanced Go concepts'),
        ('rust intro', 'Introduction to Rust');
    `)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
func TestSearchNotes_Match(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()
	ctx := context.Background()
	pattern := "%go%"
	notes, err := SearchNotes(ctx, db, pattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(notes)
}

func TestSearchNotes_NoMatch(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()
	ctx := context.Background()
	pattern := "%python%"
	notes, err := SearchNotes(ctx, db, pattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(notes)
}

func TestSearchNotes_InjectionAttempt(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()
	ctx := context.Background()
	beforeNotes, err := ListNotes(ctx, db)
	if err != nil {
		t.Fatalf("Failed to list notes: %v", err)
	}
	if len(beforeNotes) != 3 {
		t.Fatalf("Expected 3 notes before injection, got %d", len(beforeNotes))
	}
	pattern := "%'; DROP TABLE notes; --"
	notes, err := SearchNotes(ctx, db, pattern)
	if err != nil {
		t.Fatalf("SearchNotes should handle injection attempt gracefully, got error: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("Expected empty slice (no matches), got %d matches", len(notes))
	}

	if notes == nil {
		t.Fatalf("Expected non-nil empty slice, got nil")
	}

	afterNotes, err := ListNotes(ctx, db)
	if err != nil {
		t.Fatalf("Failed to list notes after injection attempt: %v", err)
	}

	if len(afterNotes) != 3 {
		t.Fatalf("Expected 3 notes after injection, got %d (table might be corrupted!)", len(afterNotes))
	}

	expected := map[int64]string{
		1: "go basics",
		2: "go advanced",
		3: "rust intro",
	}

	for _, note := range afterNotes {
		expectedTitle, exists := expected[note.ID]
		if !exists {
			t.Fatalf("Unexpected note with ID %d: %s", note.ID, note.Title)
		} else if note.Title != expectedTitle {
			t.Fatalf("Title mismatch for ID %d: expected %s, got %s",
				note.ID, expectedTitle, note.Title)
		}
		t.Log(note)
	}
}
