package sessions

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	UserAgent *string
	CreatedAt time.Time
}
