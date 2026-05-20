package http

import (
	"fmt"
	"io"
	"net/http"
)

func getWithParams() {
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/get", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	q := req.URL.Query()
	q.Set("max_price", "100")
	q.Set("min_price", "10")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/138.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
	}
	fmt.Println("URL: ", req.URL.String())
	fmt.Println("status: ", resp.StatusCode)
	fmt.Println("body: ", string(body))

}
