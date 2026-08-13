package notesrepositorygo

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	repo := Repository{
		db: db,
	}
	return &repo
}

func entityToNote(entity NoteEntity) Note {
	var body *string
	if entity.Body.Valid {
		value := entity.Body.String
		body = &value
	}

	return Note{
		entity.ID,
		entity.Title,
		body,
		entity.CreatedAt,
	}
}
