package main

import (
	"log"
	"net/http"
)

//- в main.go создать mux := http.NewServeMux()
//- зарегистрировать хэндлер на GET /health, отвечающий 200 OK и телом ok (с Content-Type: text/plain; charset=utf-8)
//- запустить через http.ListenAndServe(":8080", mux)
//- при ошибке ListenAndServe (не http.ErrServerClosed) — log.Fatal
//- проверить вручную: curl -i http://localhost:8080/health → видим 200 OK и ok
//- проверить: curl -i http://localhost:8080/missing → видим 404
//- проект кладём в один файл main.go, никаких сторонних библиотек

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	err := http.ListenAndServe(":8080", mux)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("error server:", err)
	}
}
