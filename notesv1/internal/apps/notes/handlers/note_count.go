package handlers

import (
	"net/http"
)

type CountNotesResponse struct {
	Count int64 `json:"count"`
}

func (h *Handlers) CountNotes(w http.ResponseWriter, r *http.Request) {
	result, err := h.usecases.NoteCount(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, CountNotesResponse{Count: result.Count})
}
