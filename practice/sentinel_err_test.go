package practice

import "testing"

func TestSentinelErr(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		id := -1
		HandleUserRequest(id)
	})
	t.Run("tc 2:", func(t *testing.T) {
		id := 403
		HandleUserRequest(id)
	})
	t.Run("tc 3:", func(t *testing.T) {
		id := 1000
		HandleUserRequest(id)
	})
	t.Run("tc 4:", func(t *testing.T) {
		id := 100
		HandleUserRequest(id)
	})
}
