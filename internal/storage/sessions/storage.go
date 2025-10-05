package sessions

import (
	"context"
	"errors"
	"fmt"

	"cloud-notes/internal/database/postgres"
	"cloud-notes/internal/database/redis"
	"cloud-notes/internal/logger"

	"github.com/google/uuid"
)

type storage struct {
	log logger.Logger
	pg  *postgres.Postgres
	rd  *redis.Redis
}

func New(log logger.Logger, pg *postgres.Postgres, rd *redis.Redis) Storage {
	return &storage{
		log: log,
		pg:  pg,
		rd:  rd,
	}
}

func (s *storage) scan(ctx context.Context, row postgres.Row) (*Session, error) {
	const op = "storage.sessions.scan"
	log := s.log.With(logger.String("op", op))

	session := new(Session)
	err := row.Scan(&session.ID, &session.UserID,
		&session.UserAgent, &session.CreatedAt)
	if err != nil && errors.Is(err, postgres.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *storage) Create(ctx context.Context, session *Session) error {
	const op = "storage.sessions.Create"
	log := s.log.With(logger.String("op", op))

	const sql = `INSERT INTO sessions (id, user_id, user_agent, 
                 created_at) VALUES ($1, $2, $3, $4)`

	_, err := s.pg.Exec(ctx, sql, session.ID, session.UserID,
		session.UserAgent, session.CreatedAt)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) GetByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	const op = "storage.sessions.GetByID"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM sessions WHERE id = $1`

	row := s.pg.QueryRow(ctx, sql, id)

	session, err := s.scan(ctx, row)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *storage) GetByUserID(
	ctx context.Context, userID uuid.UUID) ([]*Session, error) {
	const op = "storage.sessions.GetByUserID"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM sessions WHERE user_id = $1`

	rows, err := s.pg.Query(ctx, sql, userID)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	sessions := make([]*Session, 0)
	for rows.Next() {
		session, err := s.scan(ctx, rows)
		if err != nil {
			log.ErrorContext(ctx, "", logger.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *storage) Update(ctx context.Context, session *Session) error {
	const op = "storage.sessions.Update"
	log := s.log.With(logger.String("op", op))

	const sql = `UPDATE sessions SET user_id = $1, user_agent = $2, 
                 created_at = $3 WHERE id = $4`

	_, err := s.pg.Exec(ctx, sql, session.UserID,
		session.UserAgent, session.CreatedAt, session.ID)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "storage.sessions.Delete"
	log := s.log.With(logger.String("op", op))

	const sql = `DELETE FROM sessions WHERE id = $1`

	_, err := s.pg.Exec(ctx, sql, id)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
