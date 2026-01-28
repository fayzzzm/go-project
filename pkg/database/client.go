package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBClient wraps the pgxpool.Pool.
type DBClient struct {
	Pool *pgxpool.Pool
}

// NewDBClient creates a new database client.
func NewDBClient(ctx context.Context, connString string) (*DBClient, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Optimization: standard settings for performance
	config.MaxConns = 20
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DBClient{Pool: pool}, nil
}

// Close closes the connection pool.
func (c *DBClient) Close() {
	c.Pool.Close()
}
