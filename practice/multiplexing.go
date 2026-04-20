package practice

import (
	"fmt"
	"time"
)

func multiplexing() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for v := range 3 {
			ch1 <- fmt.Sprintf("строка %d канала ch1", v)
			time.Sleep(100 * time.Millisecond)
		}
		close(ch1)
	}()

	go func() {
		for v := range 5 {
			ch2 <- fmt.Sprintf("строка %d канала ch2", v)
			time.Sleep(150 * time.Millisecond)
		}
		close(ch2)
	}()

	counter := 1

	for ch1 != nil || ch2 != nil {
		select {
		case v, ok := <-ch1:
			if !ok {
				ch1 = nil
				continue
			}
			fmt.Printf("ch1 v:%s, serial number: %d\n", v, counter)
			counter++
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Printf("ch2 v:%s, serial number: %d\n", v, counter)
			counter++
		}
	}
}
