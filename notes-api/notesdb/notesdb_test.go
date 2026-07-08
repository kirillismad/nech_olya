package notesdb

import (
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
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE t(id INTEGER)")
	if err != nil {
		t.Fatal(err)
	}
}
