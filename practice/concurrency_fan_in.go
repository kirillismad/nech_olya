package practice

import (
	"fmt"
	"sync"
)

func merge(a, b <-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	read := func(ch <-chan string) {
		defer wg.Done()
		for val := range ch {
			out <- val
		}
	}
	wg.Add(2)
	go read(a)
	go read(b)

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func fanIn() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for i := 1; i <= 3; i++ {
			str := fmt.Sprintf("строка %d", i)
			ch1 <- str
		}
		close(ch1)
	}()

	go func() {
		for i := 4; i <= 6; i++ {
			str := fmt.Sprintf("строка %d", i)
			ch2 <- str
		}
		close(ch2)
	}()
	newCh := merge(ch1, ch2)
	for val := range newCh {
		fmt.Println(val)
	}
}
