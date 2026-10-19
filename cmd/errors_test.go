package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
	"github.com/yuyudeqiu/chronicle/internal/service"
	"gorm.io/gorm"
)

func TestTaskCommandsReturnErrors(t *testing.T) {
	if err := service.InitDB(filepath.Join(t.TempDir(), "test.db")); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		cmd  *cobra.Command
		args []string
	}{
		{getCmd, []string{"missing"}},
		{updateCmd, []string{"missing"}},
		{logCmd, []string{"missing", "work"}},
	} {
		t.Run(tc.cmd.Name(), func(t *testing.T) {
			if err := tc.cmd.RunE(tc.cmd, tc.args); !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Fatalf("error = %v, want record not found", err)
			}
		})
	}
	for _, cmd := range []*cobra.Command{summaryCmd, weeklySummaryCmd} {
		t.Run(cmd.Name(), func(t *testing.T) {
			var parseErr *time.ParseError
			if err := cmd.RunE(cmd, []string{"invalid-date"}); !errors.As(err, &parseErr) {
				t.Fatalf("error = %v, want date parse error", err)
			}
		})
	}
}

func TestImportExportReturnFileErrors(t *testing.T) {
	previousDir := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = previousDir })
	if err := exportCmd.RunE(exportCmd, nil); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("export error = %v, want missing file", err)
	}
	if err := importCmd.RunE(importCmd, []string{filepath.Join(t.TempDir(), "missing.db")}); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("import error = %v, want missing file", err)
	}
}

func TestImportReturnsConfirmationReadError(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source.db")
	if err := os.WriteFile(source, nil, 0600); err != nil {
		t.Fatal(err)
	}
	importCmd.SetIn(strings.NewReader(""))
	t.Cleanup(func() { importCmd.SetIn(nil) })
	if err := importCmd.RunE(importCmd, []string{source}); !errors.Is(err, io.EOF) {
		t.Fatalf("error = %v, want confirmation EOF", err)
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func TestPrintJSONReturnsErrors(t *testing.T) {
	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)
	var unsupported *json.UnsupportedTypeError
	if err := printJSON(cmd, make(chan int)); !errors.As(err, &unsupported) {
		t.Fatalf("error = %v, want JSON encoding error", err)
	}
	if output.Len() != 0 {
		t.Fatalf("unexpected output: %s", output.String())
	}
	writeErr := errors.New("output unavailable")
	cmd.SetOut(failingWriter{writeErr})
	if err := printJSON(cmd, map[string]string{"status": "ok"}); !errors.Is(err, writeErr) {
		t.Fatalf("error = %v, want output error", err)
	}
}

func TestExecuteReturnsErrorWithoutPrintingIt(t *testing.T) {
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"unknown-command"})
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	})
	if err := Execute(); err == nil {
		t.Fatal("expected unknown command error")
	}
	if output.Len() != 0 {
		t.Fatalf("unexpected output before entry point handles error: %s", output.String())
	}
}
