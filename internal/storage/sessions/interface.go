package sessions

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id uuid.UUID) error
}
