package notesdb

import (
	"context"
	"log"
	"testing"
)

func TestOpenInMemory_PingOk(t *testing.T) {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}

func TestOpenInMemory_IsUsable(t *testing.T) {
	db, err := OpenInMemory()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE t(id INTEGER)")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMigrate_CreatesTable(t *testing.T) {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	err = Migrate(ctx, db)
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
