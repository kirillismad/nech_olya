package dto

import "notesv1/internal/apps/notes/models"

type ListNotesQuery struct {
	Limit  int
	Offset int
}

type ListNotesResult struct {
	Items []*models.Note
}
