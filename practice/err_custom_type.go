package practice

import (
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s: %s", e.Field, e.Message)
}

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("validation error: ID: %d", e.ID)
}

func ValidateUser(name string, age int) error {
	if name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "cannot be empty",
		}
	}
	if age < 0 {
		return &ValidationError{
			Field:   "age",
			Message: "cannot be negative",
		}
	}
	return nil
}

func FindUserByID(id int) error {
	if id != 123 {
		return &NotFoundError{ID: id}
	}
	return nil
}
