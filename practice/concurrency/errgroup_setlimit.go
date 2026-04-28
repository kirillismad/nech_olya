package concurrency

import (
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func errgroupSetLimit() {
	g := new(errgroup.Group)
	g.SetLimit(3)
	totalStart := time.Now()
	for i := 1; i <= 12; i++ {
		i := i
		g.Go(func() error {
			start := time.Now()
			time.Sleep(300 * time.Millisecond)
			completion := time.Since(start)
			log.Printf("Задача %v, время старта %v, время завершения задачи %v\n", i, start, completion)
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		log.Println("ошибка", err)
	}
	totalTime := time.Since(totalStart)
	log.Println("общее время выполнения задач:", totalTime)

}
