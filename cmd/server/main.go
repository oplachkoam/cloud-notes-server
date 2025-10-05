package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud-notes/internal/config"
	"cloud-notes/internal/database/postgres"
	"cloud-notes/internal/database/redis"
	"cloud-notes/internal/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := logger.MustLoad(&cfg.Logger)
	_ = postgres.MustConnect(ctx, &cfg.Postgres)
	_ = redis.MustConnect(ctx, &cfg.Redis)

	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Second * time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.IdleTimeout),
	}

	log.InfoContext(ctx, "starting server",
		logger.String("env", cfg.Env),
		logger.String("host", cfg.Server.Host),
		logger.Int("port", cfg.Server.Port))
	_ = srv.ListenAndServe()
}
