package dto

import "notesv1/internal/apps/notes/models"

type ListNotesQuery struct{}

type ListNotesResult struct {
	Items []*models.Note
}
