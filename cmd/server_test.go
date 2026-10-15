package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yuyudeqiu/chronicle/internal/model"
)

func TestAPIRouterWithoutStaticFiles(t *testing.T) {
	router := newAPIRouter()
	for _, tc := range []struct {
		path string
		want int
	}{
		{"/api/v1/version", http.StatusOK},
		{"/", http.StatusNotFound},
		{"/assets/missing.js", http.StatusNotFound},
		{"/api/v1/missing", http.StatusNotFound},
	} {
		t.Run(tc.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if recorder.Code != tc.want {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tc.want, recorder.Body)
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Fatalf("content type = %q", got)
			}
			var response model.StandardResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			wantCode := 0
			if tc.want == http.StatusNotFound {
				wantCode = http.StatusNotFound
			}
			if response.Code != wantCode {
				t.Fatalf("response code = %d, want %d", response.Code, wantCode)
			}
		})
	}
}
