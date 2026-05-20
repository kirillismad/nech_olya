package practice

import (
	"fmt"
	"io"
	"net/http"
)

func getRequest() string{
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Println(resp.Status)
		code:=resp.Status
		return code
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
		code:=resp.Status
		return code
	}
	fmt.Println(string(body))
	code:=resp.Status
	return code
}
