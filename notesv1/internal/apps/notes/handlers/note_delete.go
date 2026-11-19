package handlers

import (
	"net/http"
	"notesv1/internal/apps/notes/usecases"

	"github.com/google/uuid"
)

func (h *Handlers) DeleteNote(w http.ResponseWriter, r *http.Request) {
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

	id, err := noteID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")

		return
	}

	_, err = h.usecases.DeleteNote(r.Context(), usecases.DeleteNoteInput{ID: id, UserID: userIDParsed})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
