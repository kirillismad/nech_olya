package http

import (
	"fmt"
	"net/http"
	"time"
)

func fetch(client *http.Client, url string) (int, error) {

	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return resp.StatusCode, nil

}

func retry(client *http.Client, url string) error {

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		statusCode, err := fetch(client, url)
		if err == nil {
			return nil
		}
		lastErr = err
		if statusCode >= 400 && statusCode < 500 {
			return err
		}
		if attempt < 3 {
			fmt.Printf("attempt %d failed: %v. retrying...\n", attempt, err)
			time.Sleep(1 * time.Second)
		}
	}
	return fmt.Errorf("all attempts failed, last error: %w", lastErr)

}

func retryRequest(url string) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	err := retry(client, url)
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	fmt.Println("request completed successfully")
}
