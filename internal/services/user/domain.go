package user

import (
	"time"

	"github.com/google/uuid"
)

type GetProfileOutput struct {
	Login     string
	FirstName string
	Timezone  string
	CreatedAt time.Time
}

type UpdateProfileInput struct {
	UserID    uuid.UUID
	FirstName string
	Timezone  string
}
