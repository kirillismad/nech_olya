package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
	"notesv1/internal/logger"

	"github.com/google/uuid"
)

type UpdateNoteRequest struct {
	Title string  `json:"title"`
	Body  *string `json:"body"`
}

func (h *Handlers) UpdateNote(w http.ResponseWriter, r *http.Request) {
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
	log := logger.FromContext(r.Context())

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

	input := usecases.UpdateNoteInput{
		ID:     id,
		Title:  req.Title,
		Body:   req.Body,
		UserID: userIDParsed,
	}

	_, err = h.usecases.UpdateNote(r.Context(), input)
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

	w.WriteHeader(http.StatusOK)
}
