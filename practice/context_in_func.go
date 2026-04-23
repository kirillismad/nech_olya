package practice

import (
	"context"
	"fmt"
	"time"
)

func startTask() {
	ctx1, cancel := context.WithCancel(context.Background())
	cancel()
	//defer cancel()
	processTask(ctx1)
}

func processTask(ctx context.Context) {
	err := ctx.Err()
	switch err {
	case context.Canceled:
		fmt.Println("контекст 1 отменен", err)
		return
	case context.DeadlineExceeded:
		fmt.Println("истек таймаут контекста 1", err)
		return
	case nil:
		fmt.Println("контекст 1 активен - продолжает работу")

	}
	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	//time.Sleep(2 * time.Second)
	finishTask(ctx2)
}

func finishTask(ctx2 context.Context) {
	err := ctx2.Err()
	switch err {
	case context.Canceled:
		fmt.Println("контекст 2 отменен", err)
		return
	case context.DeadlineExceeded:
		fmt.Println("истек таймаут контекста 2", err)
		return
	case nil:
		fmt.Println("контекст 2 активен - продолжает работу")

	}
	ctx3, cancel := context.WithTimeout(ctx2, time.Second)
	defer cancel()
	//time.Sleep(2 * time.Second)
	err = ctx3.Err()
	switch err {
	case context.Canceled:
		fmt.Println("контекст 3 отменен", err)
		return
	case context.DeadlineExceeded:
		fmt.Println("истек таймаут контекста 3", err)
		return
	case nil:
		fmt.Println("контекст 3 активен - продолжает работу")
	}
}
