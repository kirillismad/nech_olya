package practice

import (
	"errors"
	"fmt"
)

type DBError struct {
	Query string
	Err   error
}

func (d *DBError) Error() string {
	return d.Query
}

func (d *DBError) Unwrap() error {
	return d.Err
}

var ErrNotFound2 = errors.New("not found")

func queryDB(query string) error {
	if query == "" {
		return &DBError{Query: query, Err: ErrNotFound2}
	}
	return nil
}

func serviceLayer(query string) error {
	err := queryDB(query)
	if err != nil {
		return fmt.Errorf("service layer: %w", err)
	}
	return nil
}
