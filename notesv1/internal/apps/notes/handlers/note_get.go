package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
	"notesv1/internal/logger"
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
	log := logger.FromContext(r.Context())

	id, err := noteID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return
	}

	input := usecases.GetNoteInput{ID: id}

	result, err := h.usecases.GetNote(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNoteNotFound):
			log.WarnContext(
				r.Context(),
				"note not found",
				slog.Any("input", input),
			)
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
