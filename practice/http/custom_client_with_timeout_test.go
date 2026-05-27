package http

import (
	"context"
	"errors"
	"testing"
)

func TestCustomClient(t *testing.T) {
	t.Run("tc 1:success ", func(t *testing.T) {
		err := customClient()
		if err != nil {
			t.Error("there shouldn't be an error")
			return
		}
		t.Log("ok")
	})

	t.Run("tc 2: ", func(t *testing.T) {
		err := customClientWithTimeout()
		if err == nil {
			t.Error("error expected")
			return
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Error("error expected: timeout")
			return
		}
		t.Log(err)

	})
}
