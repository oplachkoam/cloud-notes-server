package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, input *RegisterInput) error
	Login(ctx context.Context, input *LoginInput) (*LoginOutput, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	ChangePassword(ctx context.Context, input *ChangePasswordInput) error
}
