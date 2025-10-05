package storage

import (
	"cloud-notes/internal/storage/sessions"
	"cloud-notes/internal/storage/users"
)

type Session = sessions.Session
type User = users.User

type Storage interface {
	Sessions() sessions.Storage
	Users() users.Storage
}
