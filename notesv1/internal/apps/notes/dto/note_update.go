package dto

import (
	"time"
	"github.com/google/uuid"
)

type UpdateNoteCommand struct {
	ID        int64
	Title     string
	Body      *string
	UpdatedAt time.Time
	UserID uuid.UUID
}

type UpdateNoteResult struct{}
