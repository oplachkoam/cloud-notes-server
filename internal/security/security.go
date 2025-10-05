package security

import (
	"context"
	"time"

	"cloud-notes/internal/config"
	"cloud-notes/internal/logger"
	"cloud-notes/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CtxKey struct{}

type security struct {
	log logger.Logger
	st  storage.Storage
	sec []byte
}

func New(log logger.Logger, st storage.Storage, cfg *config.JWT) Security {
	return &security{
		log: log,
		st:  st,
		sec: []byte(cfg.Secret),
	}
}

func (s *security) GenerateAccessToken(
	_ context.Context, claims *Claims) string {
	const op = "security.GenerateAccessToken"
	_ = s.log.With(logger.String("op", op))

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID.String(),
		"session_id": claims.SessionID.String(),
		"created_at": claims.CreatedAt.String(),
	}).SignedString(s.sec)

	return token
}

func (s *security) ParseAccessToken(
	_ context.Context, accessToken string) (*Claims, error) {
	const op = "security.ParseAccessToken"
	_ = s.log.With(logger.String("op", op))

	token, err := jwt.Parse(accessToken,
		func(token *jwt.Token) (interface{}, error) {
			return s.sec, nil
		})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	userID := uuid.MustParse(claims["user_id"].(string))
	sessionID := uuid.MustParse(claims["session_id"].(string))
	createdAt, _ := time.Parse(time.RFC3339, claims["created_at"].(string))

	return &Claims{
		UserID:    userID,
		SessionID: sessionID,
		CreatedAt: createdAt,
	}, nil
}

func GetClaims(ctx context.Context) *Claims {
	return ctx.Value(CtxKey{}).(*Claims)
}

func SetClaims(ctx context.Context, claims *Claims) context.Context {
	ctx = context.WithValue(ctx, CtxKey{}, claims)
	return ctx
}
