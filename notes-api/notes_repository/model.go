package notesrepositorygo

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64
	Title     string
	Body      *string
	CreatedAt time.Time
}

type NoteEntity struct {
	ID        int64
	Title     string
	Body      sql.NullString
	CreatedAt time.Time
}
