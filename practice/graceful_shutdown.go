package practice

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func gracefulShutdown(n int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())
	tasks := make(chan int)
	count := 0

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		<-sigChan
		cancel()
	}()

	go func() {
		defer close(tasks)
		for i := 1; i <= n; i++ {
			select {
			case <-ctx.Done():
				return
			case tasks <- i:
			}
		}
	}()

	for range n {
		wg.Add(1)
		go worker2(ctx, tasks, &wg, &mu, &count)
	}
	wg.Wait()
	fmt.Println("Кол-во обработанных задач", count)

}

func worker2(ctx context.Context, tasks <-chan int, wg *sync.WaitGroup, mu *sync.Mutex, count *int) {
	defer wg.Done()
	for {
		select {
		case val, ok := <-tasks:
			if !ok {
				return
			}
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("задача %v обработана\n", val)
			mu.Lock()
			*count++
			mu.Unlock()
		case <-ctx.Done():
			fmt.Printf("кол-во выполненных задач: %d, завершение по причине: %v\n", count, ctx.Err())
			return
		}
	}
}
