package entities

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type Note struct {
	ID        int64
	Title     string
	Body      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uuid.UUID
}

const (
	NotesTable = "notes"

	NotesIDColumn        = "id"
	NotesTitleColumn     = "title"
	NotesBodyColumn      = "body"
	NotesCreatedAtColumn = "created_at"
	NotesUpdatedAtColumn = "updated_at"
	NotesUserIDColumn    = "user_id"
)
