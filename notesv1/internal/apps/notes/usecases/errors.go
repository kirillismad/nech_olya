package usecases

import "errors"

var (
	ErrNoteNotFound  = errors.New("note not found")
	ErrValidation    = errors.New("validation error")
	ErrInvalidLimit  = errors.New("invalid limit")
	ErrInvalidOffset = errors.New("invalid offset")
)
