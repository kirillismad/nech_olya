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

type ArchiveNoteResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      *string   `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handlers) ArchiveNote(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")

		return
	}

	input := usecases.ArchiveNoteInput{ID: id}

	_, err = h.usecases.ArchiveNote(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNoteNotFound):
			log.WarnContext(
				r.Context(),
				"note not found",
				slog.Any("input", input),
			)
			writeError(w, http.StatusNotFound, "note not found")
		case errors.Is(err, usecases.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
