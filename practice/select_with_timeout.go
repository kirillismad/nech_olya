package practice

import (
	"fmt"
	"time"
)

func selectWithTimeout() {
	fastChan := make(chan string, 1)
	slowlyChan := make(chan string, 1)

	go func() {
		time.Sleep(200 * time.Millisecond)
		fastChan <- "1 строка"
		close(fastChan)
	}()

	go func() {
		time.Sleep(time.Second)
		slowlyChan <- "2 строка"
		close(slowlyChan)
	}()

	select {
	case mg := <-fastChan:
		fmt.Println(mg)
	case mg := <-slowlyChan:
		fmt.Println(mg)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("сработал timeout")
	}
}
