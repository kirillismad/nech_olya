package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Name  string `json: "name"`
	Email string `json: "email"`
}

func postWithJson() {
	user := User{
		Name:  "John",
		Email: "john@gmail.com",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error while encoding object:", err)
		return
	}
	body := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", "https://httpbin.org/post", body)
	if err != nil {
		fmt.Println("error request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
