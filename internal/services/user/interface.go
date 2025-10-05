package user

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*GetProfileOutput, error)
	UpdateProfile(ctx context.Context, input *UpdateProfileInput) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}
