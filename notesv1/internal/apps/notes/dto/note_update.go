package dto

import "time"

type UpdateNoteCommand struct {
	ID        int64
	Title     string
	Body      *string
	UpdatedAt time.Time
}

type UpdateNoteResult struct{}
