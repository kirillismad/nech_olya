package dto

import "notesv1/internal/apps/notes/models"

type ListNotesQuery struct {
	Limit  int64
	Offset int64
}

type ListNotesResult struct {
	Items []*models.Note
}
