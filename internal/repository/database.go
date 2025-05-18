package repository

import (
	"context"
	"expense_tracker/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg *config.Config) (*Database, error) {

	poolConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to parse DB URL: %w", err)
	}

	poolConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		c.TypeMap().RegisterType(&pgtype.Type{
			Name:  "date",
			Codec: pgtype.DateCodec{},
		})
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("repository: unable to create connection pool: %w", err)
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
