package notes

import (
	"database/sql"
	"net/http"
	"notesv1/internal/apps/notes/handlers"
	"notesv1/internal/apps/notes/repository"
	"notesv1/internal/apps/notes/time_provider"
	"notesv1/internal/apps/notes/usecases"
)

type Package struct {
	DB *sql.DB
}

func (p *Package) ServeMux() *http.ServeMux {
	transactionManager := repository.NewTransactionManager(p.DB)
	usecases := usecases.New(&usecases.NewArgs{
		TransactionManager: transactionManager,
		TimeProvider:       time_provider.New(),
	})
	handlers := handlers.New(&handlers.NewArgs{
		Usecases: usecases,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.ListNotes)
	mux.HandleFunc("POST /", handlers.CreateNote)
	mux.HandleFunc("GET /{id}", handlers.GetNote)
	mux.HandleFunc("PUT /{id}", handlers.UpdateNote)
	mux.HandleFunc("DELETE /{id}", handlers.DeleteNote)

	return mux
}
