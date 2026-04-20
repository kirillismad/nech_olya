package practice

import (
	"testing"
	"time"
)

func TestApplicationWaitGroupTiming(t *testing.T) {
	start := time.Now()

	applicationWaitGroup(100*time.Millisecond, time.Second)
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("Слишком долго: %v, ожидалось ~2s", elapsed)
	}
}
