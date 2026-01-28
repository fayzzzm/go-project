package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExecQueryOne executes a query and returns a single row parsed into struct T.
// It handles the query execution, checking for errors, and collecting the row.
// ExecQueryOne executes a query and returns a single row parsed into struct T.
// It handles the query execution, checking for errors, and collecting the row.
func ExecQueryOne[T any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) (*T, error) {
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, MapPgError(err)
	}
	defer rows.Close()

	val, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, MapPgError(err)
	}
	return &val, nil
}

// ExecQueryUpdate executes a query and updates the destination struct with the result.
func ExecQueryUpdate[T any](ctx context.Context, pool *pgxpool.Pool, dest *T, query string, args ...any) error {
	val, err := ExecQueryOne[T](ctx, pool, query, args...)
	if err != nil {
		return err
	}
	*dest = *val
	return nil
}

// Exec executes a query without returning any rows and maps the error.
func Exec(ctx context.Context, pool *pgxpool.Pool, query string, args ...any) error {
	_, err := pool.Exec(ctx, query, args...)
	return MapPgError(err)
}

// ExecQueryList executes a query and returns a slice of rows parsed into struct T.
// It handles the query execution, checking for errors, and collecting the rows.
func ExecQueryList[T any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) ([]T, error) {
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, MapPgError(err)
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, MapPgError(err)
	}
	return results, nil
}

// MapPgError translates PostgreSQL errors to Domain errors
func MapPgError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return domain.ErrConflict
		case "23503": // foreign_key_violation
			return fmt.Errorf("%w: %s", domain.ErrInvalidInput, pgErr.Detail)
		case "22P02": // invalid_text_representation
			return fmt.Errorf("%w: %s", domain.ErrInvalidInput, pgErr.Message)
		}
	}
	return err
}
