package cmd

import (
	"path/filepath"
	"testing"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func TestUpdateCommandUpdatesTitle(t *testing.T) {
	service.InitDB(filepath.Join(t.TempDir(), "chronicle.db"))
	task, err := service.CreateTask(model.CreateTaskReq{
		Title:    "Original title",
		Category: "test",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if err := updateCmd.Flags().Set("title", "Renamed title"); err != nil {
		t.Fatalf("set title flag: %v", err)
	}
	t.Cleanup(func() {
		_ = updateCmd.Flags().Set("title", "")
	})
	updateCmd.Run(updateCmd, []string{task.ID})

	updated, err := service.GetTask(task.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Title != title {
		t.Fatalf("title = %q, want %q", updated.Title, title)
	}
}
