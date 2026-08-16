package handlers

import (
	"errors"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
	"strconv"
	"time"
)

type GetNoteResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      *string   `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handlers) GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := noteID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return
	}

	result, err := h.usecases.GetNote(r.Context(), usecases.GetNoteInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNoteNotFound):
			writeError(w, http.StatusNotFound, "note not found")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, GetNoteResponse{
		ID:        result.Note.ID,
		Title:     result.Note.Title,
		Body:      result.Note.Body,
		CreatedAt: result.Note.CreatedAt,
		UpdatedAt: result.Note.UpdatedAt,
	})
}

func noteID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
