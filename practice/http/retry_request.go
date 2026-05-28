package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	url200 = "https://jsonplaceholder.typicode.com/posts/1"
	url500 = "https://httpbin.org/status/500"
)

type Url200Resp struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Написать программу на Go которая:
// - выполняет HTTP-запрос и проверяет три уровня ошибок: сетевая ошибка (err != nil), неуспешный статус-код (>= 400), ошибка декодирования тела
// - реализует функцию retry которая повторяет запрос до 3 раз с паузой между попытками (time.Sleep)
// - повторяет только при сетевых ошибках и статусах 5xx (серверные ошибки), не повторяет при 4xx
// - логирует номер попытки и причину повтора
// - возвращает последнюю ошибку если все попытки исчерпаны
// - написать тест

func getDoRequest(client *http.Client) func() (*http.Response, error) {
	const successfulAttempt = 3
	var attempt int = 1
	return func() (*http.Response, error) {
		defer func() {
			attempt++
		}()
		if attempt >= successfulAttempt {
			return client.Get(url200)
		}
		return client.Get(url500)
	}
}

func retry2() (Url200Resp, error) {
	var lastErr error

	doRequest := getDoRequest(http.DefaultClient)

	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, canRetry, err := func() (Url200Resp, bool, error) {
			resp, err := doRequest()
			if err != nil {
				return Url200Resp{}, true, fmt.Errorf("%w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return Url200Resp{}, false, fmt.Errorf("%d", resp.StatusCode)
			}

			if resp.StatusCode >= 500 {
				return Url200Resp{}, true, fmt.Errorf("%d", resp.StatusCode)
			}

			var urlResp Url200Resp
			err = json.NewDecoder(resp.Body).Decode(&urlResp)
			if err != nil {
				return Url200Resp{}, false, fmt.Errorf("%w", err)
			}

			return urlResp, false, nil
		}()

		if err == nil {
			return result, nil
		}

		if canRetry {
			lastErr = err
			fmt.Printf("attempt number %d, reason for retry: %s\n", attempt, err)
		} else {
			fmt.Printf("attempt number %d, reason for no retry: %s\n", attempt, err)
			return Url200Resp{}, err
		}

		if attempt < maxAttempts {
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("all attempts exhausted, last error: %s\n", lastErr)
	return Url200Resp{}, lastErr
}

func retry() (Url200Resp, error) {
	var lastErr error
	doRequest := getDoRequest(http.DefaultClient)
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := doRequest()
		if err != nil {
			lastErr = err
			fmt.Printf("attempt number %d, reason for retry: %s\n", attempt, err)
			time.Sleep(time.Second)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			fmt.Printf("attempt number %d, reason for retry: %v\n", attempt, resp.StatusCode)
			return Url200Resp{}, fmt.Errorf("%d", resp.StatusCode)
		}
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("%d", resp.StatusCode)
			fmt.Printf("attempt number %d, reason for retry: %v\n", attempt, err)
			time.Sleep(time.Second)
			continue
		}
		var Url Url200Resp
		err = json.NewDecoder(resp.Body).Decode(&Url)
		if err != nil {
			return Url200Resp{}, fmt.Errorf("%w", err)
		}
	}
	return Url200Resp{}, lastErr
}
