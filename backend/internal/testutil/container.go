//go:build integration

// Package testutil provides infrastructure for integration tests that run
// against a real PostgreSQL instance. Containers are started by testcontainers-go
// and reused across tests within the same package (one container per TestMain).
package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName  = "diarygo_test"
	testDBUser  = "diarygo"
	testDBPass  = "diarygo"
	testDBImage = "postgres:14-alpine"
)

// Container bundles the running Postgres testcontainer with a pool connected to it.
type Container struct {
	Postgres *postgres.PostgresContainer
	Pool     *pgxpool.Pool
	DSN      string
}

// StartPostgres launches a throwaway Postgres container and returns a ready-to-use
// connection pool. Callers must call Terminate when done (typically via defer in TestMain).
func StartPostgres(ctx context.Context) (*Container, error) {
	pg, err := postgres.Run(
		ctx,
		testDBImage,
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres container: %w", err)
	}

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pg.Terminate(ctx)
		return nil, fmt.Errorf("get container DSN: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		_ = pg.Terminate(ctx)
		return nil, fmt.Errorf("connect pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		_ = pg.Terminate(ctx)
		return nil, fmt.Errorf("ping pool: %w", err)
	}

	return &Container{Postgres: pg, Pool: pool, DSN: dsn}, nil
}

// Terminate closes the pool and stops the container. Safe to call multiple times.
func (c *Container) Terminate(ctx context.Context) error {
	if c == nil {
		return nil
	}
	if c.Pool != nil {
		c.Pool.Close()
	}
	if c.Postgres != nil {
		return c.Postgres.Terminate(ctx)
	}
	return nil
}
