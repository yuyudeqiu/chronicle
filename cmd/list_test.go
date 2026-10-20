package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func TestListStatusFiltersCLIAndAPI(t *testing.T) {
	if err := service.InitDB(filepath.Join(t.TempDir(), "chronicle.db")); err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"todo", "in-progress", "blocked", "done", "legacy"} {
		if err := service.DB.Create(&model.Task{ID: status, Title: status, Category: "test", Status: status}).Error; err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	if err := service.DB.Create(&model.Task{ID: "archived", Title: "archived", Category: "test", Status: "todo", ArchivedAt: &now}).Error; err != nil {
		t.Fatal(err)
	}
	router := newAPIRouter()
	previousJSON := jsonOutput
	jsonOutput = true
	t.Cleanup(func() { jsonOutput = previousJSON; listCmd.SetOut(nil) })
	for _, tc := range []struct {
		filter string
		want   []string
	}{
		{"", []string{"blocked", "in-progress", "legacy", "todo"}},
		{"todo", []string{"todo"}},
		{"in-progress", []string{"in-progress"}},
		{"blocked", []string{"blocked"}},
		{"done", []string{"done"}},
		{"todo,in-progress", []string{"in-progress", "todo"}},
		{"blocked,in-progress,todo", []string{"blocked", "in-progress", "todo"}},
		{"in-progress,blocked,todo", []string{"blocked", "in-progress", "todo"}},
		{"todo,done", []string{"done", "todo"}},
		{" todo, todo ", []string{"todo"}},
		{"done,done", []string{"done"}},
		{"not-a-status", nil},
		{"todo,invalid", nil},
		{"todo,", nil},
		{",todo", nil},
		{" ", nil},
	} {
		t.Run(tc.filter, func(t *testing.T) {
			var output bytes.Buffer
			listCmd.SetOut(&output)
			var args []string
			if tc.filter != "" {
				args = []string{tc.filter}
			}
			err := listCmd.RunE(listCmd, args)
			if tc.want == nil {
				if !errors.Is(err, service.ErrInvalidTaskStatus) || output.Len() != 0 {
					t.Fatalf("CLI error = %v, output = %q", err, output.String())
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				checkListedTaskIDs(t, output.Bytes(), tc.want)
			}
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/tasks?status="+url.QueryEscape(tc.filter), nil))
			var response struct {
				Code int             `json:"code"`
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if tc.want == nil {
				if recorder.Code != http.StatusBadRequest || response.Code != 400 {
					t.Fatalf("API status = %d, body = %s", recorder.Code, recorder.Body)
				}
			} else {
				if recorder.Code != http.StatusOK || response.Code != 0 {
					t.Fatalf("API status = %d, body = %s", recorder.Code, recorder.Body)
				}
				checkListedTaskIDs(t, response.Data, tc.want)
			}
		})
	}
	var output bytes.Buffer
	listCmd.SetOut(&output)
	if err := listCmd.RunE(listCmd, []string{"archived"}); err != nil {
		t.Fatal(err)
	}
	checkListedTaskIDs(t, output.Bytes(), []string{"archived"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/tasks/archived", nil))
	var response struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("archived API status = %d", recorder.Code)
	}
	checkListedTaskIDs(t, response.Data, []string{"archived"})
}

func checkListedTaskIDs(t *testing.T, data []byte, want []string) {
	t.Helper()
	var tasks []model.ActiveTaskResp
	if err := json.Unmarshal(data, &tasks); err != nil {
		t.Fatal(err)
	}
	ids := make([]string, 0, len(tasks))
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	sort.Strings(ids)
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("task IDs = %v, want %v", ids, want)
	}
}
