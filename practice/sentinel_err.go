package practice

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidInput     = errors.New("invalid input")
)

func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("FindUser: %w", ErrNotFound)
	}
	return "ok find user", nil
}

func PermissionDenied(id int) (string, error) {
	if id == 403 {
		return "", fmt.Errorf("PermissionDenied: %w", ErrPermissionDenied)
	}
	return " not permission denied", nil
}

func InvalidInput(id int) (string, error) {
	if id >= 1000 {
		return "", fmt.Errorf("InvalidInput: %w", ErrInvalidInput)
	}
	return "correct input", nil
}
