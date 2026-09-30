package dto

import "notesv1/internal/apps/notes/models"

type GetNoteQuery struct {
	ID int64
}

type GetNoteResult struct {
	Note *models.Note
}
