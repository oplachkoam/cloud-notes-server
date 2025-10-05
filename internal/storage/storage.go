package storage

import (
	"cloud-notes/internal/database/postgres"
	"cloud-notes/internal/database/redis"
	"cloud-notes/internal/logger"
	"cloud-notes/internal/storage/sessions"
	"cloud-notes/internal/storage/users"
)

type storage struct {
	users    users.Storage
	sessions sessions.Storage
}

func New(log logger.Logger, pg *postgres.Postgres, rd *redis.Redis) Storage {
	return &storage{
		users:    users.New(log, pg, rd),
		sessions: sessions.New(log, pg, rd),
	}
}

func (s *storage) Sessions() sessions.Storage {
	return s.sessions
}

func (s *storage) Users() users.Storage {
	return s.users
}
