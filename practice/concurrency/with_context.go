package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ConcurrencyWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}
	time.Sleep(3 * time.Second)
	cancel()
	wg.Wait()
	fmt.Println("все воркеры завершили работу")

}

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("воркер %d завершил свою работу %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("воркер %d в процессе работы...\n", id)
			time.Sleep(2 * time.Second)
		}
	}
}
