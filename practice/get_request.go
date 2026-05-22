package practice

import (
	"fmt"
	"io"
	"net/http"
)

func getRequest() (int, error) {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Println(resp.Status)
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
		return resp.StatusCode, err
	}

	fmt.Println(string(body))

	return resp.StatusCode, nil
}
