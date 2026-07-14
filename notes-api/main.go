package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

//К серверу из задачи 17 добавить операции обновления и удаления и навести порядок с кодами:
// - в store реализовать Update(id int, in Note) (Note, bool) и Delete(id int) bool
// - PUT /notes/{id}: декод JSON c DisallowUnknownFields; если заметка есть — обновить title/body, отдать 200 и обновлённую заметку; если нет — 404; пустой title — 422
// - DELETE /notes/{id}: если есть — удалить и вернуть 204 без тела; если нет — 404
// - ввести в коде два хелпера: writeJSON(w, code int, v any) и writeError(w, code int, msg string) (формат {"error":"..."}); переписать все существующие ответы через них
// - убедиться, что 405 от mux всё ещё работает и содержит правильный Allow (теперь там должно появиться PUT, DELETE для /notes/{id})
// - проверить:
// - curl -i -X PUT -d '{"title":"x","body":"y"}' http://localhost:8080/notes/1
// - curl -i -X DELETE http://localhost:8080/notes/1 → 204
// - curl -i -X DELETE http://localhost:8080/notes/1 → 404

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	mu     sync.RWMutex
	items  map[int]Note
	nextID int
}

type NoteUpdateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *Store) Create(in Note) Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	in.ID = s.nextID
	in.CreatedAt = time.Now()
	s.items[in.ID] = in
	s.nextID++
	return in
}

func (s *Store) Get(id int) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, ok := s.items[id]
	if !ok {
		return Note{}, false
	}
	return note, true
}

func (s *Store) List(limit, offset int) []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Note, 0, limit)

	for id := 1; id < s.nextID; id++ {
		note, ok := s.items[id]
		if !ok {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}

		result = append(result, note)

		if len(result) == limit {
			break
		}
	}

	return result
}

func (s *Store) Update(id int, in Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = in
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	if !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func main() {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		log.Fatal("invalid token")
	}

	store := &Store{
		items:  make(map[int]Note),
		nextID: 1,
	}

	mux := http.NewServeMux()
	registerRoutes(mux, store, token)

	handler := withRecover(withLogging(mux))
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %s", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %s", err)
		return
	}
	log.Println("shutting down…")
}

func registerRoutes(mux *http.ServeMux, store *Store, token string) {
	mux.Handle("GET /docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, "ok")
	})

	notesMux := http.NewServeMux()

	// GET localhost:8080/ -> localhost:8080/notes/
	notesMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		limitURL := r.URL.Query().Get("limit")
		if limitURL == "" {
			limitURL = "20"
		}
		limit, err := strconv.Atoi(limitURL)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to retrieve query parameters")
			return
		}
		offsetURL := r.URL.Query().Get("offset")
		if offsetURL == "" {
			offsetURL = "0"
		}
		offset, err := strconv.Atoi(offsetURL)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to retrieve query parameters")
			return
		}
		if limit < 1 || limit > 100 || offset < 0 {
			writeError(w, http.StatusBadRequest, "invalid value")
			return
		}
		notes := store.List(limit, offset)
		writeJSON(w, http.StatusOK, notes)

	})
	// POST localhost:8080/ -> localhost:8080/notes/
	notesMux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var note Note
		r.Body = http.MaxBytesReader(w, r.Body, 1<<16)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&note)
		if err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				writeError(w, http.StatusRequestEntityTooLarge, "request entity too large")
				return
			}
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if note.Title == "" {
			writeError(w, http.StatusUnprocessableEntity, "empty title")
			return
		}

		created := store.Create(note)
		w.Header().Set("Location", fmt.Sprintf("/notes/%d", created.ID))
		writeJSON(w, http.StatusCreated, created)
	})
	// GET localhost:8080/{id} -> localhost:8080/notes/{id}
	notesMux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		note, ok := store.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "not found")
			return
		}

		writeJSON(w, http.StatusOK, note)
	})

	// PUT localhost:8080/{id} -> localhost:8080/notes/{id}
	notesMux.HandleFunc("PUT /{id}", func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		note, ok := store.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		var noteRequest NoteUpdateRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err = dec.Decode(&noteRequest)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if noteRequest.Title == "" {
			writeError(w, http.StatusUnprocessableEntity, "empty title")
			return
		}

		note.Body = noteRequest.Body
		note.Title = noteRequest.Title

		store.Update(id, note)
		writeJSON(w, http.StatusOK, note)
	})

	// DELETE localhost:8080/{id} -> localhost:8080/notes/{id}
	notesMux.HandleFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		_, ok := store.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "not found")
			return
		}

		store.Delete(id)
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	// GET localhost:8080/notes
	mux.Handle("/notes/", withAuth(token, http.StripPrefix("/notes", notesMux)))

	notesMux.HandleFunc("GET /{id}/render", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		strId := r.PathValue("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		note, ok := store.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "not found")
			return
		}

		done := make(chan string, 1)
		defer close(done)

		go func() {
			time.Sleep(5 * time.Second)
			done <- "<rendered html>"
		}()
		select {
		case <-ctx.Done():
			log.Printf("render cancelled: id=%d err=%v note=%v", id, ctx.Err(), note)
			return
		case res := <-done:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(res))
		}
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("failed to encode json: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": msg})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(rec, r)

		log.Printf("method=%s path=%s status=%d duration=%v remote=%s", r.Method, r.URL.Path, rec.status, time.Since(start), r.RemoteAddr)
	})
}

func withRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("%s %s", rec, debug.Stack())
				writeError(w, http.StatusInternalServerError, "internal server error")

			}
		}()
		next.ServeHTTP(w, r)
	})
}

func withAuth(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		parts := strings.Fields(auth)
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if parts[1] != token {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		next.ServeHTTP(w, r)

	})
}
