package notesdb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func helper(t *testing.T) (context.Context, *sql.DB) {
	db := NewTestDB(t)
	ctx := context.Background()
	err := Migrate(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, `
	INSERT INTO notes (title, body)
	VALUES(?,?),(?,?)
	`,
		"title 1",
		"body 1",
		"title 2",
		"body 2")
	if err != nil {
		t.Fatal(err)
	}
	return ctx, db
}

func TestGetNoteByID_Found(t *testing.T) {
	ctx, db := helper(t)
	note, err := GetNoteByID(ctx, db, 1)
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != 1 {
		t.Fatal("incorrect id")
	}
	t.Log(note)
}

func TestGetNoteByID_NotFound(t *testing.T) {
	ctx, db := helper(t)
	_, err := GetNoteByID(ctx, db, 9999)
	if !errors.Is(err, ErrNoteNotFound) {
		t.Fatal(err)
	}
}

func TestGetNoteByID_DoesNotLeakSqlError(t *testing.T) {
	ctx, db := helper(t)
	_, err := GetNoteByID(ctx, db, 9999)
	if errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
}
