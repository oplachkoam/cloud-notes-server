package users

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
	UserStatusDeleted UserStatus = "deleted"
)

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	FirstName    string
	Timezone     string
	Status       UserStatus
	CreatedAt    time.Time
}
