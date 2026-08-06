package notesdb

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Тесты (in-memory; таблицы не нужны):
// - TestSlowCount_DeadlineExceeded — с ctx, cancel := context.WithTimeout(parent, 50*time.Millisecond) функция возвращает ошибку, и errors.Is(err, context.DeadlineExceeded) истинно
// - TestSlowCount_Cancelled — запустить SlowCount в горутине; через 10 мс вызвать cancel(); через канал получить из горутины ошибку, для которой errors.Is(err, context.Canceled) истинно
// - каждый тест завершается быстро (< 3 секунд); общий запуск пакета — с флагом -race

func TestSlowCount_DeadlineExceeded(t *testing.T) {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = SlowCount(ctx, db)
	if err == nil {
		t.Fatal("expected errors")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("expected a deadline error")
	}
}

func TestSlowCount_CancelledV2(t *testing.T) {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	defer close(errCh)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_, err := SlowCount(ctx, db)
		errCh <- err
	}()
	time.Sleep(10 * time.Millisecond)
	cancel()
	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error to cancel, got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatal("expected error context.Canceled")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("test timed out")
	}

}
