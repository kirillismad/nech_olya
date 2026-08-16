package handlers

import (
	"encoding/json"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
)

type CreateNoteRequest struct {
	Title string  `json:"title"`
	Body  *string `json:"body"`
}

type CreateNoteResponse struct {
	ID int64 `json:"id"`
}

func (h *Handlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateNoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.usecases.CreateNote(ctx, usecases.CreateNoteInput{
		Title: req.Title,
		Body:  req.Body,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	writeJSON(w, http.StatusCreated, CreateNoteResponse{ID: result.ID})

}
