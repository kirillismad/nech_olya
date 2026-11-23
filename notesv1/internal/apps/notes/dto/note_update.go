package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateNoteCommand struct {
	ID        int64
	Title     string
	Body      *string
	UserID    uuid.UUID
	Tags      []string
	UpdatedAt time.Time
}

type UpdateNoteResult struct{}
