package handlers

import (
	"errors"
	"net/http"
	"notesv1/internal/apps/notes/usecases"
	"strconv"
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

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "20"
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid limit")

		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid offset")

		return
	}

	result, err := h.usecases.ListNotes(ctx, usecases.ListNotesInput{Limit: limit, Offset: offset})
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrValidation):
			writeError(w, http.StatusBadRequest, "invalid validation")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}

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
