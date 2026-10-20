package service

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
)

func TestStatusFilterKeepsRecentHistoryLimit(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 55; i++ {
		completed := base.Add(time.Duration(i) * time.Hour)
		if err := DB.Create(&model.Task{
			ID: fmt.Sprint(i), Title: "Completed", Category: "test",
			Status: model.TaskStatusDone, ActualCompletedAt: &completed,
		}).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, filter := range []string{"done", " done,done "} {
		tasks, err := GetTasksByStatus(filter)
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 50 || tasks[0].ID != "54" || tasks[49].ID != "5" {
			t.Fatalf("filter %q: expected latest 50 tasks, got %#v", filter, tasks)
		}
	}
	tasks, err := GetTasksByStatus("todo,done")
	if err != nil || len(tasks) != 55 {
		t.Fatalf("mixed filter count = %d, error = %v, want all 55 matches", len(tasks), err)
	}
	tasks, err = GetTasksByStatus("todo")
	if err != nil || tasks == nil || len(tasks) != 0 {
		t.Fatalf("empty filter result = %#v, error = %v, want empty list", tasks, err)
	}
}
