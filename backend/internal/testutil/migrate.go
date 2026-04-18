//go:build integration

package testutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// flywayName matches Flyway files like V1__auth_usuarios.sql and extracts version + name.
var flywayName = regexp.MustCompile(`^V(\d+)__(.+)\.sql$`)

// ApplyMigrations applies the Flyway-style migrations from sourceDir to the database at dsn.
// Because golang-migrate expects the format "{version}_{name}.up.sql" but the project uses
// Flyway's "V{N}__{name}.sql", files are copied into a temp directory with the converted
// names before migrate runs. Callers typically pass "../../migrations" (relative to test pkg).
func ApplyMigrations(dsn, sourceDir string) error {
	converted, cleanup, err := prepareMigrationSource(sourceDir)
	if err != nil {
		return fmt.Errorf("prepare migration source: %w", err)
	}
	defer cleanup()

	m, err := migrate.New("file://"+converted, "pgx5://"+stripScheme(dsn))
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// prepareMigrationSource copies each V{N}__name.sql into a temp dir renamed as
// {N}_name.up.sql so golang-migrate can consume it.
func prepareMigrationSource(sourceDir string) (string, func(), error) {
	abs, err := filepath.Abs(sourceDir)
	if err != nil {
		return "", nil, err
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return "", nil, err
	}

	type pair struct{ version, src, dst string }
	var files []pair

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := flywayName.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		files = append(files, pair{
			version: m[1],
			src:     filepath.Join(abs, e.Name()),
			dst:     fmt.Sprintf("%s_%s.up.sql", m[1], m[2]),
		})
	}

	if len(files) == 0 {
		return "", nil, fmt.Errorf("no migrations found in %s", abs)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].version < files[j].version })

	tmp, err := os.MkdirTemp("", "diarygo-migrate-")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { _ = os.RemoveAll(tmp) }

	for _, f := range files {
		if err := copyFile(f.src, filepath.Join(tmp, f.dst)); err != nil {
			cleanup()
			return "", nil, err
		}
	}
	return tmp, cleanup, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// stripScheme removes a leading "postgres://" / "postgresql://" so the DSN can
// be prefixed with "pgx5://" for golang-migrate.
func stripScheme(dsn string) string {
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if len(dsn) > len(prefix) && dsn[:len(prefix)] == prefix {
			return dsn[len(prefix):]
		}
	}
	return dsn
}
