package practice

import (
	"fmt"
)

func commaOk() {
	ch := make(chan int, 3)
	ch <- 0
	ch <- 42
	ch <- 0
	close(ch)

	for {
		val, ok := <-ch
		if !ok {
			fmt.Printf("Получено из канала ch value: %d-zero value, bool: %v\n", val, ok)
			break
		}
		if val == 0 {
			fmt.Printf("Получено из канала ch value: %d-real zero, bool: %v\n", val, ok)
		} else {
			fmt.Printf("Получено из канала ch value: %d, bool: %v\n", val, ok)
		}

	}

	unbufCh := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			unbufCh <- i
		}
		close(unbufCh)
	}()
	for {
		val, ok := <-unbufCh
		if !ok {
			fmt.Println("Чтение unbufCh завершено: канал закрыт")
			break
		}
		fmt.Printf("Получено из канала unbufCh value: %d, bool: %v\n", val, ok)
	}
}
