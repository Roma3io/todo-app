package main

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"todo-app/internal/config"
	"todo-app/internal/db/postgresql"
	"todo-app/internal/http-server/handlers/tasks"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("Initializing server", slog.String("address", cfg.Address))

	storage, err := postgresql.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}
	router := chi.NewRouter()
	router.Post("/tasks", tasks.Create(storage))
	router.Get("/tasks/{id}", tasks.Get(storage))
	router.Get("/tasks", tasks.GetAll(storage))
	router.Delete("/tasks/{id}", tasks.Delete(storage))
	router.Post("/tasks/{id}", tasks.Update(storage))
	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	log.Info("Starting server", slog.String("address", cfg.Address))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
