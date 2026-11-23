package dto

import (
	"notesv1/internal/apps/notes/models"

	"github.com/google/uuid"
)

type ListNotesQuery struct {
	Limit  int64
	Offset int64
	UserID uuid.UUID
	Tag    *string
}

type ListNotesResult struct {
	Items []*models.Note
}
