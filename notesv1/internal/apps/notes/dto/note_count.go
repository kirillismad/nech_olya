package dto

import "github.com/google/uuid"

type CountNotesQuery struct {
	UserID uuid.UUID
}

type CountNotesResult struct {
	Count int64
}
