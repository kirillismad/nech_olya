package practice

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	Id int
	T  time.Duration
}

type Result struct {
	TaskID   int
	Success  bool
	Duration time.Duration
}

func MiniProject() {
	var wg sync.WaitGroup
	if len(os.Args) < 3 {
		fmt.Println("Использование: go run main.go [кол-во workers] [список длительностей задач в ms]")
		return
	}

	jobs := make(chan Task, len(os.Args)-2)
	results := make(chan Result, len(os.Args)-2)
	numWorker, _ := strconv.Atoi(os.Args[1]) //кол-во воркеров

	for i := 2; i < len(os.Args); i++ {
		duration, _ := strconv.Atoi(os.Args[i])
		jobs <- Task{Id: i - 1, T: time.Duration(duration) * time.Millisecond} //генератор задач
	}
	close(jobs)

	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// агрегатор
	successfulID := []int{}
	timeoutID := []int{}
	totalTask := 0
	totalTime := time.Duration(0)

	for result := range results {
		totalTask++
		totalTime += result.Duration
		if result.Success == true {
			successfulID = append(successfulID, result.TaskID)
		} else {
			timeoutID = append(timeoutID, result.TaskID)
		}

	}

	fmt.Printf("Общее количество задач: %d\n", totalTask)
	fmt.Printf("Количество успешных задач: %d\n", len(successfulID))
	fmt.Printf("Количество timeout: %d\n", len(timeoutID))
	fmt.Printf("Суммарное время работы: %v\n", totalTime)
	fmt.Printf("ID успешных задач: %v\n", successfulID)
	fmt.Printf("ID задач с timeout: %v\n", timeoutID)

}

func worker(jobs <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range jobs {
		done := make(chan bool)
		go func(t Task) {
			time.Sleep(t.T)
			done <- true
		}(task)

		select {
		case <-done:
			results <- Result{TaskID: task.Id, Success: true, Duration: task.T}
		case <-time.After(800 * time.Millisecond):
			results <- Result{TaskID: task.Id, Success: false, Duration: task.T}
		}
	}
}
