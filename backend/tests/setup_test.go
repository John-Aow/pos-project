package tests

import (
	"os"
	"path/filepath"
	"testing"
)

const migrationDirectory = "../infrastructure/postgres/migrations"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func MigrationDirectory(t *testing.T) string {
	t.Helper()

	absolute, err := filepath.Abs(migrationDirectory)
	if err != nil {
		t.Fatalf("resolve migration directory: %v", err)
	}

	if info, err := os.Stat(absolute); err != nil {
		t.Fatalf("migration directory is not available at %s: %v", absolute, err)
	} else if !info.IsDir() {
		t.Fatalf("migration path is not a directory: %s", absolute)
	}

	return absolute
}
