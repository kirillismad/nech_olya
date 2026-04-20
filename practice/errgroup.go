package practice

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func errgroupPractice() {
	g := new(errgroup.Group)
	start := time.Now()

	for i := 0; i < 3; i++ {
		g.Go(func() error {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			return nil
		})
	}

	for i := 3; i < 5; i++ {
		i := i
		g.Go(func() error {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			return fmt.Errorf("error task %d", i)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("done 5 task")
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
