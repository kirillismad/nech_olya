package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//К серверу из задачи 15 добавить чтение параметров пути и query:
// - зарегистрировать GET /notes/{id}, доставать id через r.PathValue("id")
// - распарсить id в int через strconv.Atoi; если не парсится — отвечать 400 и телом {"error":"invalid id"}
// - пока хранилища нет — отвечать заглушкой {"id":<id>,"title":"stub"}
// - в GET /notes прочитать query-параметры limit и offset через r.URL.Query().Get(...):
// - значения по умолчанию: limit=20, offset=0
// - валидация: limit в диапазоне 1..100, offset >= 0, иначе 400 с понятным сообщением
// - в ответе включать выбранные значения, например {"items":[],"limit":20,"offset":0}
// - проверить вручную:
// - curl -i 'http://localhost:8080/notes?limit=5'
// - curl -i 'http://localhost:8080/notes?limit=999' → 400
// - curl -i http://localhost:8080/notes/42 → 200 stub
// - curl -i http://localhost:8080/notes/abc → 400

func main() {
	mux := http.NewServeMux()
	registerRoutes(mux)
	err := http.ListenAndServe(":8080", mux)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal("error server:", err)
	}
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		limitURL := r.URL.Query().Get("limit")
		if limitURL == "" {
			limitURL = "20"
		}
		limit, err := strconv.Atoi(limitURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to retrieve query parameters"))
			return
		}
		offsetURL := r.URL.Query().Get("offset")
		if offsetURL == "" {
			offsetURL = "0"
		}
		offset, err := strconv.Atoi(offsetURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to retrieve query parameters"))
			return
		}
		if limit < 1 || limit > 100 || offset < 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid value"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"items":[],"limit":%d,"offset":%d}`, limit, offset)))
	})
	mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))
	})
	mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		strId := r.PathValue("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"id":%d,"title":"stub"}`, id)))

	})
}
