package practice

import (
	"context"
	"fmt"
	"time"
)

func contextCreationCancellation() {
	func() {
		ctx, cancel := context.WithCancel(context.Background())
		// defer cancel() тоже самое, только с принтами, для наглядности
		defer func() {
			fmt.Println("defer start")
			cancel()
			fmt.Println("defer end")
		}()

		go func() {
			fmt.Println("goroutine start")
			<-ctx.Done()
			fmt.Printf("goroutine end (%v)\n", ctx.Err())
		}()
		time.Sleep(100 * time.Millisecond)
	}()
	time.Sleep(500 * time.Millisecond)
}
