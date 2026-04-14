package practice

import (
	"fmt"
	"math"
	"sync"
)

type Result struct {
	WorkerID int
	Value    int
}

func workerPool() {
	var wg sync.WaitGroup
	jobs := make(chan int, 10)
	result := make(chan Result, 10)

	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, result, &wg)
	}

	wg.Wait()
	close(result)

	sum := 0
	for val := range result {
		fmt.Printf("Воркер: %d, Квадрат числа: %d\n", val.WorkerID, val.Value)
		sum++
	}
	fmt.Println("Общее число обработанных задач: ", sum)

}

func worker(id int, jobs <-chan int, result chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		y := math.Pow(float64(job), 2)
		result <- Result{WorkerID: id, Value: int(y)}
	}
}
