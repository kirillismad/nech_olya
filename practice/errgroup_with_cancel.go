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
		for v := range 10 {
			log.Printf("Итерация %d в горутине 1\n", v+1)
			if err := ctx.Err(); err != nil {
				log.Printf("Цикл остановлен на итерации %d: %v\n", v+1, err)
				return err
			}
			time.Sleep(50 * time.Millisecond)
		}
		return fmt.Errorf("горутина 1 вернула ошибку")
	})
	for i := 2; i < 5; i++ {
		g.Go(func() error {
			log.Printf("%d горутина начала работу\n", i)
			for v := range 10 {
				if err := ctx.Err(); err != nil {
					log.Printf("Горутина %d остановлен на итерации %d: %v, время остановки: %v\n", i, v+1, err, time.Now())
					return err
				}
				time.Sleep(50 * time.Millisecond)
			}
			return nil
		})
	}
	err := g.Wait()
	fmt.Println("Возвращенная ошибка:", err)
}
