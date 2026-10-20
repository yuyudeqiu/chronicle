package service

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"gorm.io/gorm"
)

func TestArchiveAndUnarchiveTask(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}
	task := createTestTask(t)
	if err := UpdateProgress(task.ID, model.UpdateProgressReq{LogText: "Keep this worklog"}); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name     string
		update   func(string) error
		archived bool
	}{
		{"archive", ArchiveTask, true},
		{"archive again", ArchiveTask, true},
		{"unarchive", UnarchiveTask, false},
		{"unarchive again", UnarchiveTask, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.update(task.ID); err != nil {
				t.Fatal(err)
			}
			stored, err := GetTask(task.ID)
			if err != nil {
				t.Fatal(err)
			}
			if (stored.ArchivedAt != nil) != tc.archived {
				t.Fatalf("archived_at = %v, want archived = %v", stored.ArchivedAt, tc.archived)
			}
			if stored.Status != task.Status || len(stored.Logs) != 1 || stored.Logs[0].LogText != "Keep this worklog" {
				t.Fatalf("task status or worklogs changed: %#v", stored)
			}
			if err := tc.update("missing"); !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Fatalf("missing task error = %v, want record not found", err)
			}
		})
	}
}
