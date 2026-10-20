package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func TestUpdateCommandUpdatesTitle(t *testing.T) {
	if err := service.InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatalf("initialize database: %v", err)
	}
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
	if err := updateCmd.RunE(updateCmd, []string{task.ID}); err != nil {
		t.Fatalf("update task: %v", err)
	}

	updated, err := service.GetTask(task.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Title != title {
		t.Fatalf("title = %q, want %q", updated.Title, title)
	}
}

func TestUpdateCommandCombinesStatusAndFields(t *testing.T) {
	if err := service.InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}
	task, err := service.CreateTask(model.CreateTaskReq{Title: "Original", Category: "test"})
	if err != nil {
		t.Fatal(err)
	}
	flags := map[string]string{
		"new-status": "done", "title": "Final title", "category": "release",
		"desc": "Verified", "target": "Ship", "links": "https://example.com",
		"deadline": "2026-10-20T18:00:00+08:00",
	}
	for name, value := range flags {
		if err := updateCmd.Flags().Set(name, value); err != nil {
			t.Fatal(err)
		}
	}
	previousJSON := jsonOutput
	jsonOutput = true
	var output bytes.Buffer
	updateCmd.SetOut(&output)
	t.Cleanup(func() {
		for name := range flags {
			_ = updateCmd.Flags().Set(name, "")
		}
		jsonOutput = previousJSON
		updateCmd.SetOut(nil)
	})
	if err := updateCmd.RunE(updateCmd, []string{task.ID}); err != nil {
		t.Fatal(err)
	}
	stored, err := service.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	check := func(task *model.Task) {
		t.Helper()
		if task.Title != "Final title" || task.Status != model.TaskStatusDone || task.Category != "release" ||
			task.Description != "Verified" || task.Targets != "Ship" || task.Links != "https://example.com" || task.ActualCompletedAt == nil {
			t.Fatalf("combined update did not save all fields: %#v", task)
		}
		wantDeadline, _ := time.Parse(time.RFC3339, flags["deadline"])
		if task.Deadline == nil || !task.Deadline.Equal(wantDeadline) {
			t.Fatalf("deadline = %v, want %v", task.Deadline, wantDeadline)
		}
		if len(task.Logs) != 0 {
			t.Fatalf("update unexpectedly created logs: %#v", task.Logs)
		}
	}
	check(stored)
	var response model.Task
	if err := json.Unmarshal(output.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	check(&response)
}
