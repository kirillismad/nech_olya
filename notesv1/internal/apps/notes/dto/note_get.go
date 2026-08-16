package dto

type GetNoteQuery struct {
	ID int64
}

type GetNoteResult struct {
	Note *Note
}
