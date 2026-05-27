package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func cancelRequestWithTimeout() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/delay/5", nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("timeout:", err)
			return err
		}
		fmt.Println("request error:", err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

func cancelRequest() error {

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/delay/5", nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("timeout:", err)
			return err
		}
		fmt.Println("request error:", err)
		return err
	}
	done := make(chan struct{})
	errChan := make(chan error, 1)
	go func() {
		defer close(done)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("timeout:", err)
				errChan <- err
				return
			}
			fmt.Println("request error:", err)
			errChan <- err
			return
		}
		defer resp.Body.Close()
		fmt.Println("request success:", resp.Status)
	}()
	time.Sleep(1 * time.Second)

	cancel()

	<-done
	err = <-errChan
	return err
}
