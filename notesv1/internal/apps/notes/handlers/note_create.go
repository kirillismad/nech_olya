package handlers

import (
	"encoding/json"
	"net/http"
	"notesv1/internal/apps/notes/usecases"

	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Title string  `json:"title"`
	Body  *string `json:"body"`
}

type CreateNoteResponse struct {
	ID int64 `json:"id"`
}

func (h *Handlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "X-User-ID header is required")
		return
	}

	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "X-User-ID header is not a valid UUID")
		return
	}

	ctx := r.Context()

	var req CreateNoteRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := h.usecases.CreateNote(ctx, usecases.CreateNoteInput{
		Title:  req.Title,
		Body:   req.Body,
		UserID: userIDParsed,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create note")

		return
	}

	writeJSON(w, http.StatusCreated, CreateNoteResponse{ID: result.ID})
}
