package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNoteCommand struct {
	Title     string
	Body      *string
	UserID    uuid.UUID
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateNoteResult struct {
	ID int64
}
