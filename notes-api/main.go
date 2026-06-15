package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//К серверу из задачи 16 добавить настоящее in-memory хранилище заметок и убрать заглушки:
// - объявить тип Note { ID int; Title string; Body string; CreatedAt time.Time } с json-тегами id, title, body, created_at
// - сделать тип store с полями mu sync.RWMutex, items map[int]Note, nextID int; методы Create(in Note) Note, Get(id int) (Note, bool), List(limit, offset int) []Note
// - store создаётся в main и передаётся хэндлерам (через замыкание или метод-handler)
// - POST /notes: декодировать тело через json.NewDecoder с DisallowUnknownFields(); валидировать, что title не пуст (иначе 422); ошибки декодирования — 400; в ответе — 201, заголовок Location: /notes/<id>, тело — созданная заметка JSON-ом
// - GET /notes/{id}: возвращать заметку из store; если нет — 404 {"error":"not found"}
// - GET /notes: возвращать срез List(limit, offset) как JSON-массив
// - все ответы JSON отдавать с Content-Type: application/json до WriteHeader
// - проверить полный цикл через curl: создать заметку, получить по id, получить список

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

func (s *Store) Create(in Note)Note {
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

	result := make([]Note, 0,limit)
	
	for id:=1;id<s.nextID;id++{
		note,ok:=s.items[id]
		for !ok{
			continue
		}
		if offset > 0 {
			offset--
			continue
		}

		result = append(result, note)

		if len(result) >= limit {
			break
		}
	}

	return result
}

func main() {
	store:=&Store{
		items: make(map[int]Note),
		nextID: 1,
	}
	
	mux := http.NewServeMux()
	registerRoutes(mux,store)
	err := http.ListenAndServe(":8080", mux)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal("error server:", err)
	}
}

func registerRoutes(mux *http.ServeMux,store *Store) {
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
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error":"failed to retrieve query parameters"})
			return
		}
		offsetURL := r.URL.Query().Get("offset")
		if offsetURL == "" {
			offsetURL = "0"
		}
		offset, err := strconv.Atoi(offsetURL)
		if err != nil {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error":"failed to retrieve query parameters"})
			return
		}
		if limit < 1 || limit > 100 || offset < 0 {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error":"invalid value"})
			return
		}
		w.WriteHeader(http.StatusOK)
		notes := store.List(limit, offset)
		json.NewEncoder(w).Encode(notes)
	})
	mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
		var note Note
		dec:=json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err:=dec.Decode(&note)
		 if err != nil {
            http.Error(w, "invalid json", http.StatusBadRequest)
            return
        }
		if note.Title==""{
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{"error":"empty title"})
			return
		}
        
        created := store.Create(note)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location",fmt.Sprintf("/notes/%d",created.ID))
		w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(created)

	})
	mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		strId := r.PathValue("id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error":"invalid id"})
			return
		}
        
        note, ok := store.Get(id)
        if !ok {
            w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error":"not found"})
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(note)

	})
}
