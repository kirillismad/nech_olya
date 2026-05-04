package practice

import (
	"errors"
	"testing"
)

func TestChainOfWrappers(t *testing.T) {
	var dbError *DBError
	t.Run("tc 1:", func(t *testing.T) {
		err := serviceLayer("")
		if err != nil {
			if errors.As(err, &dbError) {
				t.Log(dbError.Query)
			} else {
				t.Error("type error expected dbError")
			}
		}
	})

	t.Run("tc 2:", func(t *testing.T) {
		err := serviceLayer("4")
		if err != nil {
			if errors.As(err, &dbError) {
				t.Error("type dbError not expected")
			}
		}
		t.Log("everything is fine")
	})

	t.Run("tc 3:", func(t *testing.T) {
		err := serviceLayer("")
		if !errors.Is(err, ErrNotFound2) {
			t.Error("error expected ErrNotFound")
		} else {
			t.Log("everything is fine")
		}
	})
}
