package user

import (
	"time"
)

type GetProfileResponse struct {
	Login     string    `json:"login"`
	FirstName string    `json:"first_name"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	Timezone  string `json:"timezone"`
}
