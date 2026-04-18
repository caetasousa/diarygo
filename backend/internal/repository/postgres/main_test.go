//go:build integration

package postgres_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/caetasousa/diarygo/internal/testutil"
)

// testDB e compartilhado por todos os testes do pacote postgres_test.
// StartPostgres sobe um container uma unica vez (TestMain), aplica todas as
// migrations em backend/migrations e mantem o container ativo ate o fim da
// suite. Entre testes, usar testutil.Truncate(t, testDB) para zerar dados
// mantendo o schema.
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	c, err := testutil.StartPostgres(ctx)
	if err != nil {
		log.Fatalf("testutil.StartPostgres: %v", err)
	}

	// Resolver caminho absoluto para backend/migrations a partir deste pacote.
	// $PWD durante `go test` do pacote postgres_test =
	// .../backend/internal/repository/postgres — subir 3 diretorios.
	migrationsDir, err := filepath.Abs(filepath.Join("..", "..", "..", "migrations"))
	if err != nil {
		_ = c.Terminate(ctx)
		log.Fatalf("resolve migrations dir: %v", err)
	}

	if err := testutil.ApplyMigrations(c.DSN, migrationsDir); err != nil {
		_ = c.Terminate(ctx)
		log.Fatalf("testutil.ApplyMigrations: %v", err)
	}

	testDB = c.Pool
	code := m.Run()

	_ = c.Terminate(ctx)
	os.Exit(code)
}
