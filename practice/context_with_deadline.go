package practice

import (
	"context"
	"fmt"
	"time"
)

func contextWithDeadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	select {
	case <-time.After(3 * time.Second):
	case <-ctx.Done():
		fmt.Println("timeout expired", ctx.Err())
	}
	time.Sleep(2 * time.Second)
}
