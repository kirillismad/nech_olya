package practice

import (
	"fmt"
	"sync"
	"time"
)

func applicationWaitGroup(delay1, delay2 time.Duration) {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("Старт 1 горутины")
		time.Sleep(delay1)
		fmt.Println("Завершение 1 горутины")
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Старт 2 горутины")
		time.Sleep(delay2)
		fmt.Println("Завершение 2 горутины")
	}()
	wg.Wait()
	fmt.Println("Все задачи завершены")
}
