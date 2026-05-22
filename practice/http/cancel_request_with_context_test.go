package http

import (
	"context"
	"errors"
	"testing"
)

func TestCancelRequest(t *testing.T) {
	t.Run("tc 1: context deadline exceeded", func(t *testing.T) {
		err := cancelRequestWithTimeout()
		if err == nil {
			t.Error("error expected")
			return
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Error("error expected: context deadline exceeded")
			return
		}
		t.Log("error:", err)
	})

	t.Run("tc 2: context canceled", func(t *testing.T) {
		err := cancelRequest()
		if err == nil {
			t.Error("error expected: context canceled")
			return
		}
		if !errors.Is(err, context.Canceled) {
			t.Error("error expected: context canceled")
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			t.Error("error expected: context canceled")
			return
		}
		t.Log("error: ",err)
	})

}
