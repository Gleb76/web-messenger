package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Gleb76/web-messenger/internal/config"
	"github.com/Gleb76/web-messenger/internal/httpx"
	"github.com/go-chi/chi/v5"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	http   *http.Server
	logger *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) *Server {
	r := chi.NewRouter()
	r.Use(httpx.RequestIDMiddleware(logger))

	r.Get("/healthz", healthzJSON)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", pingJSON)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			httpx.WriteError(w, http.StatusNotFound, "not_found", "route not found")
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "route not found")
	})

	return &Server{
		http: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Handler() http.Handler { return s.http.Handler }

func (s *Server) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)
	go func() {
		s.logger.Info("server starting", "addr", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
	case err := <-serverErr:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("graceful shutdown failed", "err", err)
		_ = s.http.Close()
		return err
	}

	s.logger.Info("server stopped")
	return nil
}

func healthzJSON(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func pingJSON(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"pong": true})
}
