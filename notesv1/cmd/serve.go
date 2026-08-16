package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"notesv1/internal/apps/notes"
	"notesv1/internal/config"
	"notesv1/internal/logger"
	"notesv1/internal/middlewares"
	pkgconfig "notesv1/pkg/config"
	"time"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		log := logger.New()
		slog.SetDefault(log)

		cfg, err := pkgconfig.LoadConfig[config.Config]("NOTES_V1")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		db, err := OpenDB(cfg)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		if err := db.PingContext(ctx); err != nil {
			return fmt.Errorf("failed to ping database: %w", err)
		}

		notesPackage := &notes.Package{
			DB: db,
		}

		mux := http.NewServeMux()
		mux.HandleFunc("GET /alive", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		mux.Handle("/notes/", http.StripPrefix("/notes", notesPackage.ServeMux()))

		middlewares := []middlewares.Middleware{
			middlewares.NewLoggerMiddleware(log),
			middlewares.NewRecoverMiddleware(),
		}

		var handler http.Handler = mux
		for _, m := range middlewares {
			handler = m(handler)
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
			Handler: handler,
		}

		errCh := make(chan error, 1)
		go func() {
			slog.Info(
				"starting server",
				slog.String("addr", srv.Addr),
			)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("failed to start server: %w", err)
				return
			}
		}()

		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			slog.Info("shutting down server")
			if err := srv.Shutdown(shutdownCtx); err != nil {
				return err
			}
			return nil
		case err := <-errCh:
			return err
		}
	},
}
