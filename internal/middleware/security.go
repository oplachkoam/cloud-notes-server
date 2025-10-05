package middleware

import (
	"errors"
	"net/http"
	"strings"

	"cloud-notes/internal/logger"
	"cloud-notes/internal/render"
	"cloud-notes/internal/security"
	"cloud-notes/internal/storage"
)

var (
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthScheme = errors.New("invalid auth scheme")
	ErrInvalidToken      = errors.New("invalid auth token")
	ErrSessionExpired    = errors.New("session expired")
)

func Security(log logger.Logger, st storage.Storage,
	sec security.Security) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.Security"
			_ = log.With(logger.String("op", op))
			ctx := r.Context()

			header := r.Header.Get("Authorization")
			if header == "" {
				render.Error(w, http.StatusUnauthorized, ErrEmptyAuthHeader)
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			if token == header {
				render.Error(w, http.StatusUnauthorized, ErrInvalidAuthScheme)
				return
			}

			claims, err := sec.ParseAccessToken(ctx, token)
			if err != nil {
				render.Error(w, http.StatusUnauthorized, ErrInvalidToken)
				return
			}

			session, err := st.Sessions().GetByID(ctx, claims.SessionID)
			if err != nil {
				render.ServerError(w, http.StatusInternalServerError)
				return
			}

			if session == nil {
				render.Error(w, http.StatusUnauthorized, ErrSessionExpired)
				return
			}

			ctx = security.SetClaims(ctx, claims)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
