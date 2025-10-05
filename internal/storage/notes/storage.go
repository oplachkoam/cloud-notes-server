package notes

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

func (s *storage) scan(ctx context.Context, row postgres.Row) (*Note, error) {
	const op = "storage.notes.scan"
	log := s.log.With(logger.String("op", op))

	note := new(Note)
	err := row.Scan(
		&note.ID, &note.UserID, &note.Title, &note.Text,
		&note.Pinned, &note.UpdatedAt, &note.CreatedAt)
	if err != nil && errors.Is(err, postgres.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return note, nil
}

func (s *storage) Create(ctx context.Context, note *Note) error {
	const op = "storage.notes.Create"
	log := s.log.With(logger.String("op", op))

	const sql = `INSERT INTO notes (id, user_id, title, text, pinned, 
                 updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.pg.Exec(
		ctx, sql, note.ID, note.UserID, note.Title, note.Text,
		note.Pinned, note.UpdatedAt, note.CreatedAt)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) GetByID(ctx context.Context, id uuid.UUID) (*Note, error) {
	const op = "storage.notes.GetByID"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM notes WHERE id = $1`

	row := s.pg.QueryRow(ctx, sql, id)

	note, err := s.scan(ctx, row)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return note, nil
}

func (s *storage) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Note, error) {
	const op = "storage.notes.GetByUserID"
	log := s.log.With(logger.String("op", op))

	const sql = `SELECT * FROM notes WHERE user_id = $1`

	rows, err := s.pg.Query(ctx, sql, userID)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	notes := make([]*Note, 0)
	for rows.Next() {
		note, err := s.scan(ctx, rows)
		if err != nil {
			log.ErrorContext(ctx, "", logger.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *storage) Update(ctx context.Context, note *Note) error {
	const op = "storage.notes.Update"
	log := s.log.With(logger.String("op", op))

	const sql = `UPDATE notes SET user_id = $1, title = $2, 
                 text = $3, pinned = $4, updated_at = $6, 
                 created_at = $7 WHERE id = $8`

	_, err := s.pg.Exec(
		ctx, sql, note.UserID, note.Title, note.Text,
		note.Pinned, note.UpdatedAt, note.CreatedAt, note.ID)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "storage.notes.Delete"
	log := s.log.With(logger.String("op", op))

	const sql = `DELETE FROM notes WHERE id = $1`

	_, err := s.pg.Exec(ctx, sql, id)
	if err != nil {
		log.ErrorContext(ctx, "", logger.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
