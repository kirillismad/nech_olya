package practice

import "testing"

func TestGrecefulShutdown(t *testing.T) {
	gracefulShutdown(100)
}
