package users

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	List(ctx context.Context, limit, offset *uint64) ([]*User, error)
	Count(ctx context.Context) (uint64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
