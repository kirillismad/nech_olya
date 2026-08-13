package notesrepositorygo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func newTestRepo(t *testing.T) *Repository {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})
	repo := New(db)
	ctx := context.Background()
	err = repo.Migrate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func TestCreate_OK(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	id, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	if id == 0 {
		t.Fatal("invalid id")
	}
	t.Log(id)
}

func TestCreate_EmptyTitle(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "",
		Body:  nil,
	}
	ctx := context.Background()
	_, err := repo.Create(ctx, note)
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func TestGetByID_Found(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	id, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	receivedNote, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receivedNote)
}

func TestGetByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	_, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.GetByID(ctx, 999)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
		t.Log(err)
	}
}

func TestList_Empty(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	ctx := context.Background()
	notes, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 0 {
		t.Fatal("expected len 0")
	}
	if notes == nil {
		t.Fatal("the not should not be nil")
	}
	t.Log("ok")
}

func TestList_ReturnsAll(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	ctx := context.Background()
	note := Note{
		Title: "title",
		Body:  nil,
	}
	_, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	notes, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 3 {
		t.Fatal("expected length 3")
	}
	t.Log(notes)
}

func TestUpdate_OK(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	id, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	body := "body"
	note2 := Note{
		Title: "Title 2",
		Body:  &body,
	}
	err = repo.Update(ctx, id, note2)
	if err != nil {
		t.Fatal(err)
	}
	note, err = repo.GetByID(ctx, id)
	if note.ID != id || note.Title != note2.Title || *note.Body != *note2.Body {
		t.Fatal("the record has not changed")
	}
	t.Log("ok")
}

func TestUpdate_NotFound(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	err := repo.Update(ctx, 999, note)

	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestDelete_Ok(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	id, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.GetByID(ctx, id)

	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestDelete_NotFound(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx := context.Background()
	_, err := repo.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.Delete(ctx, 999)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestBulkCreate_Success(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	ctx := context.Background()
	notes := []Note{{Title: "title 1", Body: nil}, {Title: "title 2", Body: nil}, {Title: "title 3", Body: nil}}
	ids, err := repo.BulkCreate(ctx, notes)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 3 {
		t.Fatal("expected length 3")
	}
	notes2, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes2) != 3 {
		t.Fatal("expected length 3")
	}
	t.Log("ok")
}

func TestBulkCreate_AllOrNothing(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	ctx := context.Background()
	notes := []Note{{Title: "", Body: nil}, {Title: "title 2", Body: nil}, {Title: "title 3", Body: nil}}
	_, err := repo.BulkCreate(ctx, notes)
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatal(err)
	}

	notes2, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes2) != 0 {
		t.Fatal("expected length 3")
	}
	t.Log("ok")
}

func TestContextCancellation(t *testing.T) {
	t.Parallel()
	repo := newTestRepo(t)
	note := Note{
		Title: "title",
		Body:  nil,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := repo.Create(ctx, note)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
		t.Log("ok 1")
	}

	_, err = repo.List(ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
		t.Log("ok 2")
	}

}
