package context

import (
	"context"
	"fmt"
	"time"
)

func ContextCreationCancellation() {
	func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			<-ctx.Done()
			fmt.Println(ctx.Err())

		}()
		time.Sleep(2 * time.Second)
	}()
	time.Sleep(5 * time.Second)
}
