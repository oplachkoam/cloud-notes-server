package storage

import (
	"cloud-notes/internal/storage/sessions"
	"cloud-notes/internal/storage/users"
)

const (
	UserStatusPending = users.StatusPending
	UserStatusActive  = users.StatusActive
	UserStatusBlocked = users.StatusBlocked
	UserStatusDeleted = users.StatusDeleted
)

type Session = sessions.Session
type User = users.User

type Storage interface {
	Sessions() sessions.Storage
	Users() users.Storage
}
