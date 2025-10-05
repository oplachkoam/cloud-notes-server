package notes

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     *string
	Text      *string
	Pinned    bool
	UpdatedAt *time.Time
	CreatedAt time.Time
}
