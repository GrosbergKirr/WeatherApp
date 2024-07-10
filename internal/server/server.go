package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	srv http.Server
}

func NewServer(cfg *internal.Config, router chi.Router) *Server {
	return &Server{srv: http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  cfg.ServerTimeout,
		WriteTimeout: cfg.ServerTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout},
	}
}
func (s *Server) ServerRun(log *slog.Logger, cfg *internal.Config) {
	log.Info("starting server", slog.String("address", cfg.ServerAddress))
	err := s.srv.ListenAndServe()
	if err != nil {
		return
	}
}

func (s *Server) ServerStop(ctx context.Context, log *slog.Logger) {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Error("Server stop error", slog.Any("err", err))
	}
	log.Info("Server successfully stopped")

}
