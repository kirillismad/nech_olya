package practice

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func ErrGroupWithCancel() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		log.Println("1 горутина начала работу")
		if err := ctx.Err(); err != nil {
			log.Printf("Горутина 1 остановлена: %v\n", err)
			return err
		}
		time.Sleep(500 * time.Millisecond)
		return fmt.Errorf("горутина 1 вернула ошибку")
	})
	for i := 2; i < 5; i++ {
		g.Go(func() error {
			log.Printf("%d горутина начала работу\n", i)
			for {
				select {
				case <-ctx.Done():
					log.Printf("Горутина %d остановлена, время остановки: %v\n", i, time.Now())
					if err := ctx.Err(); err != nil {
						return err
					}
				case <-time.After(100 * time.Millisecond):

				}
			}
			return nil
		})
	}
	err := g.Wait()
	log.Println("Возвращенная ошибка:", err)
}
