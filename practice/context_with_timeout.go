package practice

import (
	"context"
	"fmt"
	"time"
)

func contextWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	select {
	case <-time.After(3 * time.Second):
	case <-ctx.Done():
		fmt.Println("timeout expired\n", ctx.Err())
	}
	time.Sleep(2 * time.Second)

}
