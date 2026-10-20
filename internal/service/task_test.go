package service

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
)

func TestUpdateProgressRejectsInvalidStatus(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatalf("initialize database: %v", err)
	}
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
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatalf("initialize database: %v", err)
	}
	task := createTestTask(t)

	_, err := UpdateTask(task.ID, model.UpdateTaskReq{Status: "doing", Title: "Should not be saved"})
	if !errors.Is(err, ErrInvalidTaskStatus) {
		t.Fatalf("UpdateTask() error = %v, want invalid task status error", err)
	}
	stored, err := GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Title != task.Title || stored.Status != task.Status {
		t.Fatalf("invalid update changed task: %#v", stored)
	}
}

func TestUpdateTaskCompletionTime(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}
	task := createTestTask(t)
	completed, err := UpdateTask(task.ID, model.UpdateTaskReq{Status: model.TaskStatusDone})
	if err != nil {
		t.Fatal(err)
	}
	if completed.Status != model.TaskStatusDone || completed.ActualCompletedAt == nil {
		t.Fatalf("task not completed: %#v", completed)
	}
	// Use an older timestamp so rewriting completion time cannot pass unnoticed.
	original := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	if err := DB.Model(task).Update("actual_completed_at", original).Error; err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"", model.TaskStatusDone} {
		updated, err := UpdateTask(task.ID, model.UpdateTaskReq{Status: status, Title: "Still done"})
		if err != nil {
			t.Fatal(err)
		}
		if updated.ActualCompletedAt == nil || !updated.ActualCompletedAt.Equal(original) {
			t.Fatalf("completion time = %v, want %v", updated.ActualCompletedAt, original)
		}
	}
	reopened, err := UpdateTask(task.ID, model.UpdateTaskReq{Status: model.TaskStatusTodo, Title: "Reopened"})
	if err != nil {
		t.Fatal(err)
	}
	if reopened.Status != model.TaskStatusTodo || reopened.ActualCompletedAt != nil || reopened.Title != "Reopened" {
		t.Fatalf("task not reopened correctly: %#v", reopened)
	}
	completedAgain, err := UpdateTask(task.ID, model.UpdateTaskReq{Status: model.TaskStatusDone})
	if err != nil {
		t.Fatal(err)
	}
	if completedAgain.ActualCompletedAt == nil || !completedAgain.ActualCompletedAt.After(original) {
		t.Fatalf("new completion time = %v, want later than %v", completedAgain.ActualCompletedAt, original)
	}
}

func TestGetActiveTasksIncludesExistingUnknownStatuses(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatalf("initialize database: %v", err)
	}
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

func TestUpdateProgressCompletesTaskWithWorklog(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatalf("initialize database: %v", err)
	}
	task := createTestTask(t)
	if err := UpdateProgress(task.ID, model.UpdateProgressReq{
		LogText:   "Completed and verified",
		NewStatus: model.TaskStatusDone,
	}); err != nil {
		t.Fatalf("UpdateProgress() error = %v", err)
	}

	stored, err := GetTask(task.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if stored.Status != model.TaskStatusDone || stored.ActualCompletedAt == nil {
		t.Fatalf("task status = %q, completed at = %v, want done with completion time", stored.Status, stored.ActualCompletedAt)
	}
	if len(stored.Logs) != 1 || stored.Logs[0].LogText != "Completed and verified" {
		t.Fatalf("task logs = %#v, want completion worklog", stored.Logs)
	}
}
