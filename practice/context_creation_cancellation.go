package practice

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func contextCreationCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		<-ctx.Done()
		fmt.Println(ctx.Err())

	}()
	wg.Wait()

}
