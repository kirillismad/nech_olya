package dto

import "github.com/google/uuid"

type DeleteNoteCommand struct {
	ID     int64
	UserID uuid.UUID
}

type DeleteNoteResult struct{}
