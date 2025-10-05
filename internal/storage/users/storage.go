package users

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

func (s *storage) scan(ctx context.Context, row postgres.Row) (*User, error) {
	const op = "storage.users.scan"
	log := s.log.With(logger.String("op", op))

	user := new(User)
	err := row.Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.FirstName,
		&user.Timezone, &user.Status, &user.CreatedAt)
	if err != nil && errors.Is(err, postgres.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *storage) Create(ctx context.Context, user *User) error {
	const op = "storage.users.Create"
	log := s.log.With(logger.String("op", op))

	const sql = `INSERT INTO users (id, login, password_hash, 
                 first_name, timezone, status, created_at) 
                 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.pg.Exec(
		ctx, sql, user.ID, user.Login, user.PasswordHash,
		user.FirstName, user.Timezone, user.Status, user.CreatedAt)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	const op = "storage.users.GetByID"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM users WHERE id = $1`

	row := s.pg.QueryRow(ctx, sql, id)

	user, err := s.scan(ctx, row)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *storage) GetByLogin(ctx context.Context, login string) (*User, error) {
	const op = "storage.users.GetByLogin"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM users WHERE login = $1`

	row := s.pg.QueryRow(ctx, sql, login)

	user, err := s.scan(ctx, row)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *storage) List(
	ctx context.Context, limit, offset *uint64) ([]*User, error) {
	const op = "storage.users.List"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM users LIMIT $1 OFFSET $2`

	rows, err := s.pg.Query(ctx, sql, limit, offset)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user, err := s.scan(ctx, rows)
		if err != nil {
			log.ErrorContext(ctx, "", logger.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *storage) Count(ctx context.Context) (uint64, error) {
	const op = "storage.users.Count"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT COUNT(*) FROM users`

	row := s.pg.QueryRow(ctx, sql)
	count := new(uint64)
	err := row.Scan(&count)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return *count, nil
}

func (s *storage) Update(ctx context.Context, user *User) error {
	const op = "storage.users.Update"
	log := s.log.With(logger.String("op", op))

	const sql = `UPDATE users SET login = $1, password_hash = $2, 
                 first_name = $3, timezone = $4, status = $6, 
                 created_at = $7 WHERE id = $8`

	_, err := s.pg.Exec(
		ctx, sql, user.Login, user.PasswordHash, user.FirstName,
		user.Timezone, user.Status, user.CreatedAt, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "storage.users.Delete"
	log := s.log.With(logger.String("op", op))

	const sql = `DELETE FROM users WHERE id = $1`

	_, err := s.pg.Exec(ctx, sql, id)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
