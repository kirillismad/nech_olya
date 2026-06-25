package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//Написать CLI-утилиту на Go для работы с REST API (например, JSONPlaceholder — https://jsonplaceholder.typicode.com):
// - утилита принимает команду через os.Args: list, get , create
// - команда list — GET-запрос к /posts, вывод первых 5 постов (id и title)
// - команда get — GET-запрос к /posts/{id}, вывод полного поста (title, body, userId)
// - команда create — POST-запрос к /posts с JSON-телом (title и body из аргументов* или захардкоженные)- использовать кастомный http.Client с таймаутом 10 секунд
// - все запросы создавать через http.NewRequestWithContext с контекстом - JSON-ответы декодировать в Go-структуры - реализовать retry для сетевых ошибок (максимум 3 попытки) - корректно обрабатывать все ошибки и выводить понятные сообщения пользователю

type Post struct {
	UserID int64  `json:"userId"`
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type APIClient struct {
	BaseURL string
	Client  *http.Client
	Num int
}

type ListPostQuery struct {
	Page     int64
	PageSize int64
}

type GetPostQuery struct {
	ID string
}

type CreatePostCommand struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userID"`
}

func (api APIClient) ListPost(ctx context.Context, query ListPostQuery) ([]Post, error) {
	if query.Page < 1 {

		return nil, fmt.Errorf("page must be greater than 0")
	}
	if query.PageSize < 1 {
		return nil, fmt.Errorf("page_size must be greater than 0")
	}
	start := (query.Page - 1) * query.PageSize

	u, err := url.Parse(api.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/posts"
	q := u.Query()
	q.Set("_start", strconv.FormatInt(start, 10))
	q.Set("_limit", strconv.FormatInt(query.PageSize, 10))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.retry(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var posts []Post

	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}

	return posts, nil

}

func (api APIClient) GetPost(ctx context.Context, query GetPostQuery) (Post, error) {
	u, err := url.Parse(api.BaseURL)
	if err != nil {
		return Post{}, err
	}
	u.Path = "/posts/" + query.ID

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Post{}, err
	}

	resp, err := api.retry(req)
	if err != nil {
		return Post{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Post{}, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var post Post

	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		return Post{}, err
	}

	return post, nil

}

func (api APIClient) CreatePost(ctx context.Context, command CreatePostCommand) (Post, error) {

	u, err := url.Parse(api.BaseURL)
	if err != nil {
		return Post{}, err
	}
	u.Path = "/posts"

	jsonBody, err := json.Marshal(command)
	if err != nil {
		return Post{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return Post{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := api.retry(req)
	if err != nil {
		return Post{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Post{}, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var createdPost Post

	err = json.NewDecoder(resp.Body).Decode(&createdPost)
	if err != nil {
		return Post{}, err
	}
	return createdPost, nil

}

func (api APIClient) retry(req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response

	for attempt := 1; attempt <= api.Num; attempt++ {
		resp, err = api.Client.Do(req)
		if err != nil {
			fmt.Println("network error:", err)
			continue
		}
		if resp.StatusCode >= 500 {
			fmt.Println("invalid status code:", resp.StatusCode)
			err = fmt.Errorf("invalid status code: %d", resp.StatusCode)
			continue
		}
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("request timeout:", err)
			continue
		}
		return resp, nil
	}
	return nil, err
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("request canceled")
		os.Exit(1)
	}()

	strValue := os.Getenv("RETRY_COUNT")
	if strValue == "" {
		strValue = "0"
	}

	num, err := strconv.Atoi(strValue)
	if err != nil {
		fmt.Println("type conversion error: ", err)
		return
	}
	if num == 0 {
		num = 3
	}

	if len(os.Args) < 2 {
		fmt.Println("command required")
		return
	}
	ctx := context.Background()

	api := APIClient{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Num: num,
	}

	command := os.Args[1]

	switch command {
	case "list":
		if len(os.Args) < 4 {
			fmt.Println("usage: list <page> <page_size>")
			return
		}

		page, err := strconv.ParseInt(os.Args[2], 10, 64)

		if err != nil {
			fmt.Println("page must be number")
			return
		}
		pageSize, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println("page_size must be number")
			return
		}
		posts, err := api.ListPost(ctx, ListPostQuery{
			Page:     page,
			PageSize: pageSize,
		})

		if err != nil {
			fmt.Println("list error:", err)
			return
		}

		for _, post := range posts {
			fmt.Printf("ID: %d, Title: %s\n", post.ID, post.Title)
		}

	case "get":
		if len(os.Args) < 3 {

			fmt.Println("id required")
			return
		}
		post, err := api.GetPost(ctx, GetPostQuery{
			ID: os.Args[2],
		})
		if err != nil {
			fmt.Println("get error:", err)
			return
		}
		fmt.Printf("UserID: %d\nID: %d\nTitle: %s\nBody: %s\n",
			post.UserID,
			post.ID,
			post.Title,
			post.Body,
		)
	case "create":
		if len(os.Args) < 5 {

			fmt.Println("title and body required")
			return
		}
		id, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		post, err := api.CreatePost(ctx, CreatePostCommand{
			Title:  os.Args[2],
			Body:   os.Args[3],
			UserID: id,
		})
		if err != nil {
			fmt.Println("create error:", err)
			return
		}
		fmt.Printf("Created post:\nID: %d\nTitle: %s\nBody: %s\nUserID: %d\n",
			post.ID,
			post.Title,
			post.Body,
			post.UserID,
		)
	default:
		fmt.Println("unknown command:", command)
	}
}
