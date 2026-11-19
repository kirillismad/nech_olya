package handlers

import (
	"net/http"
	"notesv1/internal/apps/notes/usecases"

	"github.com/google/uuid"
)

type CountNotesResponse struct {
	Count int64 `json:"count"`
}

func (h *Handlers) CountNotes(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.usecases.NoteCount(r.Context(), usecases.NoteCountInput{UserID: userIDParsed})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, CountNotesResponse{Count: result.Count})
}
