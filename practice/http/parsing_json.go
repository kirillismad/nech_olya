package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func parsingJSONUnmarshal() {
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		fmt.Println("error request:", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:", err)
		return
	}
	var post Post
	err = json.Unmarshal(bodyBytes, &post)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(post.UserID)
	fmt.Println(post.ID)
	fmt.Println(post.Title)
	fmt.Println(post.Body)
}

func parsingJSONDecoder() {
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		fmt.Println("error request:", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var post Post
	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(post.UserID)
	fmt.Println(post.ID)
	fmt.Println(post.Title)
	fmt.Println(post.Body)
}
