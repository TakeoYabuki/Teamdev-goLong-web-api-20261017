package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORSAddsAllowOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)

	WithCORS(next).ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
	if rec.Code != http.StatusTeapot {
		t.Errorf("中の処理がよばれていません: status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}
