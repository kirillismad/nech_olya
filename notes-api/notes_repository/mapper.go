package notesrepositorygo

func entityToNote(entity NoteEntity) Note {
	var body *string
	if entity.Body.Valid {
		value := entity.Body.String
		body = &value
	}

	return Note{
		ID:        entity.ID,
		Title:     entity.Title,
		Body:      body,
		CreatedAt: entity.CreatedAt,
	}
}
