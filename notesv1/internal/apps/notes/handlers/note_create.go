package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"notesv1/internal/apps/notes/usecases"

	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Title string   `json:"title"`
	Body  *string  `json:"body"`
	Tags  []string `json:"tags"`
}

type CreateNoteResponse struct {
	ID   int64    `json:"id"`
	Tags []string `json:"tags"`
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
		Tags:   req.Tags,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNoteNotFound):
			writeError(w, http.StatusNotFound, "note not found")
			return
		case errors.Is(err, usecases.ErrValidation):
			writeError(w, http.StatusBadRequest, "invalid validation")
			return
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusCreated, CreateNoteResponse{ID: result.ID, Tags: result.Tags})
}
