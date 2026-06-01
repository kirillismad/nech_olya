package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
}

type ListPostQuery struct {
	Page     int64
	PageSize int64
}

type GetPostQuery struct {
	ID string
}

type CreatePostCommand struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (api APIClient) ListPost(ctx context.Context, query ListPostQuery) ([]Post, error) {
	if query.Page < 1 {

		return nil, fmt.Errorf("page must be greater than 0")
	}
	if query.PageSize < 1 {
		return nil, fmt.Errorf("page_size must be greater than 0")
	}
	start := (query.Page - 1) * query.PageSize
	url := api.BaseURL + "/posts" +
		"?_start=" + strconv.FormatInt(start, 10) +
		"&_limit=" + strconv.FormatInt(query.PageSize, 10)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := retry(api.Client, req)
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

	url := api.BaseURL + "/posts/" + query.ID

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Post{}, err
	}

	resp, err := retry(api.Client, req)
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

	url := api.BaseURL + "/posts"

	jsonBody, err := json.Marshal(command)
	if err != nil {
		return Post{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return Post{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := retry(api.Client, req)
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

func retry(client *http.Client, req *http.Request) (resp *http.Response, err error) {
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("network error:", err)
			continue
		}
		return resp, nil
	}
	return nil, err
}

func main() {
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
		if len(os.Args) < 4 {

			fmt.Println("title and body required")
			return
		}
		post, err := api.CreatePost(ctx, CreatePostCommand{
			Title: os.Args[2],
			Body:  os.Args[3],
		})
		if err != nil {
			fmt.Println("create error:", err)
			return
		}
		fmt.Printf("Created post:\nID: %d\nTitle: %s\nBody: %s\n",
			post.ID,
			post.Title,
			post.Body,
		)
	default:
		fmt.Println("unknown command:", command)
	}
}
