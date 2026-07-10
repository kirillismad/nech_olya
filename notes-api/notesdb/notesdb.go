package notesdb

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func OpenInMemory() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open: %w", err)
	}

	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY AUTOINCREMENT, 
title TEXT NOT NULL, body TEXT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return fmt.Errorf("frequest could not be completed: %w", err)
	}
	return nil
}

func InsertNote(ctx context.Context,db *sql.DB,title,body string)(int64,error){
	result,err:=db.ExecContext(ctx,`INSERT INTO notes(title,body) VALUES(?,?)`,title,body)
	if err != nil {
		return 0,fmt.Errorf("frequest could not be completed: %w", err)
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		return 0,fmt.Errorf("failed to retrieve id: %w", err)
	}

	return id,nil
}
