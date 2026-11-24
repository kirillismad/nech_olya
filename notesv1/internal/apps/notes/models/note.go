package models

import (
	"time"
	"github.com/google/uuid"
)

type Note struct {
	ID        int64
	Title     string
	Body      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uuid.UUID
}
