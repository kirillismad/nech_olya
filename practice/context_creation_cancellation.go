package practice

import (
	"context"
	"fmt"
	"time"
)

func contextCreationCancellation() {
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
