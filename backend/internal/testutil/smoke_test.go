//go:build integration

package testutil_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/testutil"
	"github.com/caetasousa/diarygo/internal/testutil/factories"
)

var testPool *pgxpool.Pool

// TestMain demonstrates the canonical lifecycle: one Postgres container per
// package, migrations applied once, pool shared across tests. Copy this pattern
// into new integration suites (e.g. internal/repository/postgres/*_test.go).
func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := testutil.StartPostgres(ctx)
	if err != nil {
		log.Fatalf("start postgres: %v", err)
	}

	if err := testutil.ApplyMigrations(container.DSN, "../../migrations"); err != nil {
		_ = container.Terminate(ctx)
		log.Fatalf("apply migrations: %v", err)
	}
	testPool = container.Pool

	code := m.Run()
	_ = container.Terminate(ctx)
	os.Exit(code)
}

func TestSmoke_FactoriesAndTruncate(t *testing.T) {
	// Creates a Cliente end-to-end through factories.
	c := factories.NewCliente(t, testPool)
	if c.ID.String() == "" {
		t.Fatal("cliente id vazio")
	}

	var count int
	if err := testPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM clientes").Scan(&count); err != nil {
		t.Fatalf("count clientes: %v", err)
	}
	if count != 1 {
		t.Fatalf("esperava 1 cliente, got %d", count)
	}

	// Truncate wipes the data but keeps the schema.
	testutil.Truncate(t, testPool)

	if err := testPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM clientes").Scan(&count); err != nil {
		t.Fatalf("count clientes pos-truncate: %v", err)
	}
	if count != 0 {
		t.Fatalf("esperava 0 clientes pos-truncate, got %d", count)
	}

	// A second factory call after truncate must still work (no unique conflicts).
	_ = factories.NewCliente(t, testPool)
}
