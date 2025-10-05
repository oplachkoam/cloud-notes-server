package auth

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)

type RegisterInput struct {
	Login     string
	Password  string
	FirstName string
	Timezone  string
}

type LoginInput struct {
	Login     string
	Password  string
	UserAgent *string
}

type LoginOutput struct {
	AccessToken string
}

type ChangePasswordInput struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}
