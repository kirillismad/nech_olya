package notesrepositorygo

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	repo := Repository{
		db: db,
	}
	return &repo
}
