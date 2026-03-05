package repository

import (
	"context"
	"expense_tracker/internal/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// The Database struct provides a connection pool to PostgreSQL.
type Database struct {
	Pool   *pgxpool.Pool
	Logger *logrus.Logger
}

// NewDB initializes a new PostgreSQL connection pool.
func NewDB(ctx context.Context, cfg *config.Config) (*Database, error) {
	log := cfg.Logger.WithFields(logrus.Fields{
		"package": "repository",
		"method":  "NewDB",
	})
	log.Info("initializing database connection pool")

	poolConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		log.WithError(err).Error("failed to parse DB URL")
		return nil, fmt.Errorf("repository: failed to parse DB URL: %w", err)
	}

	poolConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		log.Debug("registering custom types for new connection")
		c.TypeMap().RegisterType(&pgtype.Type{
			Name:  "date",
			Codec: pgtype.DateCodec{},
		})
		return nil
	}

	log.WithField("max_connections", poolConfig.MaxConns).Debug("creating connection pool")
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.WithError(err).Error("unable to create connection pool")
		return nil, fmt.Errorf("repository: unable to create connection pool: %w", err)
	}

	return &Database{
		Pool:   pool,
		Logger: cfg.Logger,
	}, nil
}

// Close closes all connections in the pool.
func (db *Database) Close() {
	log := db.Logger.WithFields(logrus.Fields{
		"package": "repository",
		"method":  "Close",
	})
	log.Info("closing database connection pool")

	start := time.Now()
	statsBefore := db.Pool.Stat()

	db.Pool.Close()
	log.WithFields(logrus.Fields{
		"connections_closed": statsBefore.TotalConns(),
		"duration_ms":        time.Since(start).Milliseconds(),
	}).Info("database connection pool closed")
}
