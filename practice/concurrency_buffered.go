package practice

import (
	"fmt"
)

func bufferedChan() {
	jobs := make(chan int, 3)

	for i := 1; i <= 3; i++ {
		jobs <- i
		fmt.Printf("len(jobs): %v, cap(jobs): %v\n", len(jobs), cap(jobs))
	}

	close(jobs)

	for job := range jobs {
		fmt.Printf("job: %v\n", job)
	}
}
