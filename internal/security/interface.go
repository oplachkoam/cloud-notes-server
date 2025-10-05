package security

import (
	"context"
)

type Security interface {
	GenerateAccessToken(ctx context.Context, claims *Claims) string
	ParseAccessToken(ctx context.Context, accessToken string) (*Claims, error)
}
