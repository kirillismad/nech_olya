package dto

import (
	"notesv1/internal/apps/notes/models"
)

type RestoreNoteQuery struct {
	ID int64
}

type RestoreNoteResult struct {
	Note *models.Note
}