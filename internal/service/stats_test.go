package service

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
)

func TestSummaryStats(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}

	// Empty databases still return seven daily entries and zero totals.
	empty, err := GetStatsSummary()
	if err != nil || empty.TotalTasks != 0 || empty.CompletionRate != 0 || len(empty.WeeklyStats) != 7 {
		t.Fatalf("empty stats = %+v, error = %v", empty, err)
	}
	emptyWeek, err := GetWeeklySummary("")
	if err != nil || emptyWeek.Stats != (model.WeeklyStats{}) || len(emptyWeek.Activities) != 0 {
		t.Fatalf("empty week = %+v, error = %v", emptyWeek, err)
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	createdAt := today.Add(time.Hour)
	statuses := []string{"todo", "in-progress", "blocked", "done", "done", "legacy", "done"}
	for i, status := range statuses {
		task := model.Task{
			ID: fmt.Sprint(i), Title: "task", Category: fmt.Sprint(i % 2),
			Status: status, CreatedAt: createdAt, UpdatedAt: createdAt,
		}
		if status == model.TaskStatusDone {
			task.ActualCompletedAt = &createdAt
		}
		if i == 6 {
			task.ArchivedAt = &createdAt
		}
		if err := DB.Create(&task).Error; err != nil {
			t.Fatal(err)
		}
		// Multiple logs must not count the same task multiple times.
		for j := 0; j < 2; j++ {
			log := model.TaskLog{
				ID: fmt.Sprintf("%d-%d", i, j), TaskID: task.ID,
				LogText: "work", CreatedAt: createdAt,
			}
			if err := DB.Create(&log).Error; err != nil {
				t.Fatal(err)
			}
		}
	}

	stats, err := GetStatsSummary()
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalTasks != 6 || stats.CompletedTasks != 2 || stats.TodoTasks != 1 ||
		stats.InProgressTasks != 1 || stats.BlockedTasks != 1 || stats.CompletionRate != 2.0/6 {
		t.Fatalf("stats = %+v", stats)
	}
	if !reflect.DeepEqual(stats.ByCategory, map[string]int{"0": 3, "1": 3}) {
		t.Fatalf("category counts = %v", stats.ByCategory)
	}
	for i, daily := range stats.WeeklyStats {
		want := model.DailyStats{Date: today.AddDate(0, 0, i-6).Format("2006-01-02")}
		if i == 6 {
			// Archive filtering applies to totals, but not historical daily counts.
			want.Created, want.Completed = 7, 3
		}
		if daily != want {
			t.Fatalf("daily stats = %+v, want %+v", daily, want)
		}
	}

	// A task without logs and a log for a missing task must not inflate reports.
	if err := DB.Create(&model.Task{
		ID: "idle", Title: "no activity", Category: "test", Status: model.TaskStatusTodo,
		CreatedAt: createdAt, UpdatedAt: createdAt,
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Create(&model.TaskLog{
		ID: "orphan", TaskID: "missing", LogText: "orphan log", CreatedAt: createdAt,
	}).Error; err != nil {
		t.Fatal(err)
	}
	daily, err := GetDailySummary("")
	if err != nil {
		t.Fatal(err)
	}
	if len(daily.Activities) != 7 {
		t.Fatalf("daily activities = %d, want 7", len(daily.Activities))
	}
	for _, activity := range daily.Activities {
		if activity.TaskID == "" || activity.TaskTitle != "task" || len(activity.TodayLogs) != 2 {
			t.Fatalf("daily activity = %+v", activity)
		}
	}

	week, err := GetWeeklySummary("")
	if err != nil {
		t.Fatal(err)
	}
	wantWeek := model.WeeklyStats{
		TotalTasks: 7, CompletedTasks: 3, TodoTasks: 1, InProgressTasks: 1, BlockedTasks: 1,
	}
	if week.Stats != wantWeek || len(week.Activities) != 7 {
		t.Fatalf("weekly stats = %+v, activities = %d", week.Stats, len(week.Activities))
	}
}

func TestGetStatsSummaryReturnsDatabaseErrors(t *testing.T) {
	if err := InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}
	db, err := DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if stats, err := GetStatsSummary(); err == nil || stats != nil {
		t.Fatalf("stats = %+v, error = %v; want database error", stats, err)
	}
}
