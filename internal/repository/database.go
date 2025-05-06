package repository

import (
	"context"
	"expense_tracker/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg *config.Config) (*Database, error) {
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("repository: unable to create connection pool: %w", err)
	}
	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
