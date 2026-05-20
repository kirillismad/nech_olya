package errs

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
			return
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected error: %v, got: %v", ErrNotFound, err)
			return
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		id := 1
		_, err := FindUser(id)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}
