//go:build integration

package testutil

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AllTables lists every table currently created by the migrations, in any order.
// Keep this list in sync with backend/migrations/V*.sql. When a new migration
// adds a table, add it here too; Truncate will include it automatically.
var AllTables = []string{
	// V1
	"usuarios",
	// V2
	"clientes",
	"profissionais",
	"enderecos",
	"documentos",
	"referencias",
	"regioes",
	"regioes_atuacao",
	"disponibilidades",
}

// Truncate wipes every table listed in AllTables. Uses a single TRUNCATE ...
// CASCADE RESTART IDENTITY statement so FK order is irrelevant and sequences reset.
// Call this between tests (not between suites) to keep the schema but clear data.
func Truncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmt := fmt.Sprintf(
		"TRUNCATE TABLE %s RESTART IDENTITY CASCADE",
		strings.Join(AllTables, ", "),
	)
	if _, err := pool.Exec(ctx, stmt); err != nil {
		t.Fatalf("truncate failed: %v", err)
	}
}
