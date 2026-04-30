package main

import (
	// Контекст для отмены работы всех горутин.
	"context"
	// Форматированный вывод в консоль.
	"fmt"
	// Доступ к сигналам ОС.
	"os"
	// Подписка на системные сигналы.
	"os/signal"
	// Примитивы синхронизации (WaitGroup).
	"sync"
	// Атомарные операции для потокобезопасных счетчиков.
	"sync/atomic"
	// Константы сигналов POSIX (SIGTERM, SIGINT).
	"syscall"
	// Работа со временем и таймерами.
	"time"
)

const (
	// Количество воркеров.
	N = 2
)

// Result описывает результат обработки одного задания.
type Result struct {
	// Исходное значение (в примере — метка времени в миллисекундах).
	Value int64
	// Признак четности значения.
	IsEven bool
}

// String задает человекочитаемое представление Result для fmt.Println.
func (r Result) String() string {
	return fmt.Sprintf("Value: %d, IsEven: %t", r.Value, r.IsEven)
}

func main() {
	// Контекст завершается при получении сигналов остановки.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// Освобождаем ресурсы, связанные с контекстом.
	defer cancel()

	// Канал входящих заданий для воркеров.
	jobs := make(chan int64)
	// Канал результатов обработки.
	result := make(chan Result)

	// published: сколько заданий отправлено в jobs.
	var published int64
	// handled: сколько заданий обработано воркерами.
	var handled int64
	// received: сколько результатов прочитано в main.
	var received int64

	// Publisher: периодически генерирует задание и отправляет в jobs.
	go func() {
		defer func() {
			fmt.Println("Publisher stopped")
		}()
		// Закрываем jobs при остановке publisher, чтобы воркеры корректно завершились.
		defer close(jobs)
		for {
			select {
			// Останавливаемся по отмене контекста.
			case <-ctx.Done():
				return
			// Пауза перед публикацией следующего задания.
			case <-time.After(51 * time.Millisecond):
				select {
				// Повторная проверка отмены перед отправкой.
				case <-ctx.Done():
					return
				// В качестве задания отправляем текущее время в миллисекундах.
				case jobs <- time.Now().UnixMilli():
					published++
				}
			}
		}
	}()

	// WaitGroup ждет завершения всех воркеров.
	var wg sync.WaitGroup
	// Запускаем N воркеров.
	for i := range N {
		wg.Go(func() {
			defer func() {
				fmt.Printf("Worker %d stopped\n", i)
			}()
			for {
				select {
				// Завершаемся по сигналу отмены.
				case <-ctx.Done():
					return
				// Читаем очередное задание из jobs.
				case job, ok := <-jobs:
					// Если канал закрыт, работы больше нет.
					if !ok {
						return
					}
					// Имитация времени обработки задания.
					select {
					case <-ctx.Done():
						return
					case <-time.After(101 * time.Millisecond):
						// Отправка результата с дополнительной проверкой отмены.
						select {
						case <-ctx.Done():
							return
						case result <- Result{Value: job, IsEven: job%2 == 0}:
							atomic.AddInt64(&handled, 1)
						}
					}
				}
			}
		})
	}
	// Закрываем канал результатов после остановки всех воркеров.
	go func() {
		wg.Wait()
		close(result)
	}()

	// Читаем результаты, пока канал result не будет закрыт.
	for r := range result {
		fmt.Println("Received:", r)
		received++
	}
	// Печатаем итоговую статистику обработки.
	fmt.Printf("\nPublished: %d, Handled: %d, Received: %d\n", published, handled, received)
}
