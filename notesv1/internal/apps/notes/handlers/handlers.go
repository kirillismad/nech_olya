package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
)

type Handlers struct {
	usecases *usecases.Usecases
}

type NewArgs struct {
	Usecases *usecases.Usecases
}

func New(args *NewArgs) *Handlers {
	return &Handlers{
		usecases: args.Usecases,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("failed to encode json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(map[string]string{"error": msg})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
	}
}
