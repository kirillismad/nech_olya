package dto

import (
	"notesv1/internal/apps/notes/models"

	"github.com/google/uuid"
)

type GetNoteQuery struct {
	ID     int64
	UserID uuid.UUID
}

type GetNoteResult struct {
	Note *models.Note
}
