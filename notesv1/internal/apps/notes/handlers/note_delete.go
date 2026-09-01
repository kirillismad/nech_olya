package handlers

import (
	"net/http"
	"notesv1/internal/apps/notes/usecases"
)

func (h *Handlers) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := noteID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")

		return
	}

	_, err = h.usecases.DeleteNote(r.Context(), usecases.DeleteNoteInput{ID: id})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
