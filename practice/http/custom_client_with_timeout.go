package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func customClient()error {
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	urls := []string{
		"https://httpbin.org/delay/ip",
		"https://httpbin.org/get",
		"https://httpbin.org",
	}

	for _, url := range urls {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("request error:", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("request error:", err)
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		fmt.Println("successful request:", string(body))
	}
	return nil
}

func customClientWithTimeout()error {
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
	resp, err := client.Get("https://httpbin.org/delay/10")
	if err != nil {
		if os.IsTimeout(err) {
			fmt.Println("timeout:", err)
			return err
		} else {
			fmt.Println("request error:", err)
			return err
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("successful request:", string(body))
	return nil
}
