package notesrepositorygo

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
