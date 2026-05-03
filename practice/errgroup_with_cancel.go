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
	for range 3 {
		g.Go(func() error {
			log.Println("Горутина начала работу")
			for {
				select {
				case <-ctx.Done():
					log.Printf("Горутина остановлена, время остановки: %v\n", time.Now())
					if err := ctx.Err(); err != nil {
						return err
					}
				case <-time.After(100 * time.Millisecond):

				}
			}
		})
	}
	err := g.Wait()
	log.Println("Возвращенная ошибка:", err)
}
