package service

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/yuyudeqiu/chronicle/internal/model"
)

func TestUpdateProgressRejectsInvalidStatus(t *testing.T) {
	InitDB(filepath.Join(t.TempDir(), "chronicle.db"))
	task := createTestTask(t)

	err := UpdateProgress(task.ID, model.UpdateProgressReq{NewStatus: "in_progress"})
	if !errors.Is(err, ErrInvalidTaskStatus) {
		t.Fatalf("UpdateProgress() error = %v, want invalid task status error", err)
	}

	stored, err := GetTask(task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if stored.Status != model.TaskStatusTodo {
		t.Fatalf("task status = %q, want %q", stored.Status, model.TaskStatusTodo)
	}
}

func TestUpdateTaskRejectsInvalidStatus(t *testing.T) {
	InitDB(filepath.Join(t.TempDir(), "chronicle.db"))
	task := createTestTask(t)

	_, err := UpdateTask(task.ID, model.UpdateTaskReq{Status: "doing"})
	if !errors.Is(err, ErrInvalidTaskStatus) {
		t.Fatalf("UpdateTask() error = %v, want invalid task status error", err)
	}
}

func TestGetActiveTasksIncludesExistingUnknownStatuses(t *testing.T) {
	InitDB(filepath.Join(t.TempDir(), "chronicle.db"))
	task := createTestTask(t)
	if err := DB.Model(task).Update("status", "doing").Error; err != nil {
		t.Fatalf("seed unknown status: %v", err)
	}

	tasks, err := GetActiveTasks()
	if err != nil {
		t.Fatalf("GetActiveTasks() error = %v", err)
	}
	if len(tasks) != 1 || tasks[0].ID != task.ID || tasks[0].Status != "doing" {
		t.Fatalf("GetActiveTasks() = %#v, want task with existing unknown status", tasks)
	}
}

func createTestTask(t *testing.T) *model.Task {
	t.Helper()
	task, err := CreateTask(model.CreateTaskReq{Title: "test task", Category: "test"})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	return task
}
