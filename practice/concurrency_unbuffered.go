package practice

import (
	"fmt"
	"math"
	"sync"
)

func unbufferedChan() {
	worker := func(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
		defer wg.Done()
		for job := range jobs {
			result := math.Pow(float64(job), 2)
			results <- int(result)
		}

	}
	var wg sync.WaitGroup
	jobs := make(chan int)
	results := make(chan int)

	wg.Add(1)
	go worker(jobs, results, &wg)

	go func() {
		for i := 1; i < 4; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
