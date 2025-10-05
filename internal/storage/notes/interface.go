package notes

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, note *Note) error
	GetByID(ctx context.Context, id uuid.UUID) (*Note, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Note, error)
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id uuid.UUID) error
}
