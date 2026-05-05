package practice

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("FindUser: %w", ErrNotFound)
	}
	return "ok find user", nil
}
