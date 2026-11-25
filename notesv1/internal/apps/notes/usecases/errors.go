package usecases

import "errors"

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrValidation   = errors.New("validation error")
)
