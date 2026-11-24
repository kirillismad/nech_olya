package dto

import (
	"time"
	"github.com/google/uuid"
)

type CreateNoteCommand struct {
	Title     string
	Body      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uuid.UUID
}

type CreateNoteResult struct {
	ID int64
}
