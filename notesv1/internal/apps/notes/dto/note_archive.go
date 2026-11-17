package dto

import (
	"time"
	"notesv1/internal/apps/notes/models"
)

type ArchiveNoteQuery struct {
	ID int64
	ArchivedAt *time.Time
}

type ArchiveNoteResult struct {
	Note *models.Note
}