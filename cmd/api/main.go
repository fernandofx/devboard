package main

import (
	"net/http"
	"os"

	"github.com/fernandofx/devboard/internal/handler"
	"github.com/fernandofx/devboard/internal/logger"
	"github.com/fernandofx/devboard/internal/middlewares"
	"github.com/fernandofx/devboard/internal/server"
)

func main() {
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	logger := logger.New(logger.ProductionConfig())

	srv := server.New(":3003", logger)

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	srv.RegisterRoutes("GET /health", handler.NewHealthHandler())
	srv.RegisterRoutes("GET /panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic de prueba")
	}))

	if err := srv.Start(); err != nil {
		logger.Error("Error Fatal", "error", err)
		os.Exit(1)
	}

}
