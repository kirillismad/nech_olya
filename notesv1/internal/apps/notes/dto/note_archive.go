package dto

import (
	"notesv1/internal/apps/notes/models"
	"time"
)

type ArchiveNoteQuery struct {
	ID         int64
	ArchivedAt *time.Time
}

type ArchiveNoteResult struct {
	Note *models.Note
}
