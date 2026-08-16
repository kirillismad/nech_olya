package handlers

import (
	"net/http"
	"notesv1/internal/apps/notes/usecases"
	"time"
)

type ListNotesResponseItems struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      *string   `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListNotesResponse struct {
	Items []ListNotesResponseItems `json:"items"`
}

func (h *Handlers) ListNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.usecases.ListNotes(ctx, usecases.ListNotesInput{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ListNotesResponse{
		Items: make([]ListNotesResponseItems, 0, len(result.Items)),
	}
	for _, note := range result.Items {
		response.Items = append(response.Items, ListNotesResponseItems{
			ID:        note.ID,
			Title:     note.Title,
			Body:      note.Body,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, response)
}
