package dto

import "time"

type CreateNoteCommand struct {
	Title     string
	Body      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateNoteResult struct {
	ID int64
}
