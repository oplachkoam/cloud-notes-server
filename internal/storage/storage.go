package storage

import (
	"cloud-notes/internal/database/postgres"
	"cloud-notes/internal/database/redis"
	"cloud-notes/internal/logger"
	"cloud-notes/internal/storage/notes"
	"cloud-notes/internal/storage/sessions"
	"cloud-notes/internal/storage/users"
)

type storage struct {
	notes    notes.Storage
	users    users.Storage
	sessions sessions.Storage
}

func New(log logger.Logger, pg *postgres.Postgres, rd *redis.Redis) Storage {
	return &storage{
		notes:    notes.New(log, pg, rd),
		users:    users.New(log, pg, rd),
		sessions: sessions.New(log, pg, rd),
	}
}

func (s *storage) Notes() notes.Storage {
	return s.notes
}

func (s *storage) Sessions() sessions.Storage {
	return s.sessions
}

func (s *storage) Users() users.Storage {
	return s.users
}
