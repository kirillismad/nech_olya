package practice

import (
	"errors"
	"testing"
)

func TestSentinelErr(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		id := -1
		_, err := FindUser(id)
		if err == nil {
			t.Errorf("error expected: %v", ErrNotFound)
		}
		if errors.Is(err, ErrNotFound) {
			t.Log(err)
		}

	})
	t.Run("tc 2:", func(t *testing.T) {
		id := 1
		_, err := FindUser(id)
		if err == nil {
			t.Log("no error")
		}
		if errors.Is(err, ErrNotFound) {
			t.Error(err)
		}
	})
}
