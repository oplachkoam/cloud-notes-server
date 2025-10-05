package postgres

import (
	"context"
	"fmt"

	"cloud-notes/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres = pgxpool.Pool
type Tx = pgx.Tx
type Batch = pgx.Batch
type Row = pgx.Row
type Rows = pgx.Rows
type Command = pgconn.CommandTag

var ErrNoRows = pgx.ErrNoRows

func Connect(ctx context.Context, cfg *config.Postgres) (*Postgres, error) {
	p, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = p.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return p, nil
}

func MustConnect(ctx context.Context, cfg *config.Postgres) *Postgres {
	p, err := Connect(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return p
}
