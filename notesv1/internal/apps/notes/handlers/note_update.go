package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
)

type UpdateNoteRequest struct {
	Title string  `json:"title"`
	Body  *string `json:"body"`
}

func (h *Handlers) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := noteID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, err = h.usecases.UpdateNote(r.Context(), usecases.UpdateNoteInput{
		ID:    id,
		Title: req.Title,
		Body:  req.Body,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNoteNotFound):
			writeError(w, http.StatusNotFound, "note not found")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
