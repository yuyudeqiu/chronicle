package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDBReturnsErrors(t *testing.T) {
	t.Run("connection", func(t *testing.T) {
		if err := InitDB(filepath.Join(t.TempDir(), "missing", "test.db")); err == nil {
			t.Fatal("expected database connection error")
		}
	})
	t.Run("migration", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "invalid.db")
		if err := os.WriteFile(path, []byte("not a SQLite database"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := InitDB(path); err == nil {
			t.Fatal("expected invalid database error")
		}
	})
}
