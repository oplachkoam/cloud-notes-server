package users

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	StatusPending UserStatus = "pending"
	StatusActive  UserStatus = "active"
	StatusBlocked UserStatus = "blocked"
	StatusDeleted UserStatus = "deleted"
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
