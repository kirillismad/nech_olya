package main

import (
	"errors"
	"log"
	"net/http"
)

//через Go 1.22 method-паттерны на http.ServeMux:
// - GET /notes — отвечать 200 и пустым JSON-массивом [] (с Content-Type: application/json)
// - POST /notes — отвечать 201 и телом {"status":"created"} (заглушка, тело запроса пока не парсить)
// - GET /health из задачи 14 остаётся работать
// - проверить, что неверный метод даёт 405 с заголовком Allow (это делает mux автоматически):
// - curl -i -X DELETE http://localhost:8080/notes → 405, в ответе есть Allow: GET, POST
// - проверить, что неизвестный путь даёт 404:
// - curl -i http://localhost:8080/unknown → 404
// - вынести регистрацию маршрутов в функцию registerRoutes(mux *http.ServeMux)

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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	})
	mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))
	})
}
