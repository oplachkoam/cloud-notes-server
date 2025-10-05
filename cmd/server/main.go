package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud-notes/internal/config"
	"cloud-notes/internal/database/postgres"
	"cloud-notes/internal/database/redis"
	authHandler "cloud-notes/internal/handlers/auth"
	userHandler "cloud-notes/internal/handlers/user"
	"cloud-notes/internal/logger"
	"cloud-notes/internal/middleware"
	"cloud-notes/internal/security"
	authService "cloud-notes/internal/services/auth"
	userService "cloud-notes/internal/services/user"
	"cloud-notes/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := logger.MustLoad(&cfg.Logger)
	pg := postgres.MustConnect(ctx, &cfg.Postgres)
	rd := redis.MustConnect(ctx, &cfg.Redis)

	st := storage.New(log, pg, rd)
	sec := security.New(log, st, &cfg.JWT)

	authSrv := authService.New(log, st, sec)
	userSrv := userService.New(log, st)

	auth := authHandler.New(log, authSrv)
	user := userHandler.New(log, userSrv)

	r := chi.NewRouter()
	r.Use(middleware.Logging(log))
	r.Route("/api", func(r chi.Router) {
		r.With(middleware.Security(log, st, sec)).Group(func(r chi.Router) {
			r.Post("/auth/logout", auth.Logout)
			r.Post("/auth/change-password", auth.ChangePassword)
			r.Route("/user", func(r chi.Router) {
				r.Get("/profile", user.GetProfile)
				r.Put("/profile", user.UpdateProfile)
				r.Delete("/profile", user.DeleteProfile)
			})
		})
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", auth.Register)
			r.Post("/auth/login", auth.Login)
		})
	})

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
