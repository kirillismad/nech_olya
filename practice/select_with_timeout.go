package practice

import (
	"fmt"
	"sync"
	"time"
)

func selectWithTimeout() {
	var wg sync.WaitGroup
	fastChan := make(chan string, 1)
	slowlyChan := make(chan string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		fastChan <- "1 строка"
		close(fastChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		slowlyChan <- "2 строка"
		close(slowlyChan)
	}()

	wg.Wait()
	select {
	case mg := <-fastChan:
		fmt.Println(mg)
	case mg := <-slowlyChan:
		fmt.Println(mg)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("сработал timeout")
	}
}
