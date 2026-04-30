package concurrency

import "fmt"

func ProducerConsumer() {
	values := make(chan int)
	result := make(chan int)
	go producer(values)
	go consumer(values, result)
	fmt.Println(<-result)
}

func producer(out chan<- int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

func consumer(in <-chan int, result chan<- int) {
	sum := 0
	for val := range in {
		sum += val
	}
	result <- sum
	close(result)
}
