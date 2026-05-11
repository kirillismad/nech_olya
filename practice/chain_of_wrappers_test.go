package practice

import (
	"errors"
	"testing"
)

func TestChainOfWrappers(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		err := serviceLayer("")
		if err != nil {
			var dbErr *DBError
			if errors.As(err, &dbErr) {
				t.Log(dbErr.Query)
			} else {
				t.Error("type error expected dbError")
				return
			}
		}
	})

	t.Run("tc 2:", func(t *testing.T) {
		err := serviceLayer("4")
		if err != nil {
			var dbErr *DBError
			if errors.As(err, &dbErr) {
				t.Error("type dbError not expected")
				return
			}
		}
		t.Log("everything is fine")
	})

	t.Run("tc 3:", func(t *testing.T) {
		err := serviceLayer("")
		if !errors.Is(err, ErrNotFound2) {
			t.Error("error expected ErrNotFound")
			return
		} else {
			t.Log("everything is fine")
		}
	})
}
