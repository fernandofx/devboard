// Package server creates and start server
package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server struct
type Server struct {
	httpserver *http.Server
	mux        *http.ServeMux
	logger     *slog.Logger
}

// New creates a new server
func New(addr string, logger *slog.Logger) *Server {

	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second, //tiempo que dura una conexión
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		httpserver: httpServer,
		mux:        mux,
		logger:     logger,
	}
}

// Start server
func (s *Server) Start() error {

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(quit)

	servErr := make(chan error, 1)

	go func() {
		s.logger.Info("Servidor Iniciado", slog.String("addr", s.httpserver.Addr))
		if err := s.httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			servErr <- err
		}

	}()

	select {
	case err := <-servErr:
		return err
	case sig := <-quit:
		s.logger.Info("Señal recibida", slog.String("Signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := s.httpserver.Shutdown(ctx); err != nil {
		s.logger.Info("Error de Shutdown", slog.Any("Error", err))
		return err
	}

	s.logger.Info("Servidor detenido OK")

	return nil

}

// RegisterRoutes riutes
func (s *Server) RegisterRoutes(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// ServeHTTP Server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Use implementar Middlewares
func (s *Server) Use(middleware func(http.Handler) http.Handler) {
	s.httpserver.Handler = middleware(s.httpserver.Handler)
}
