package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestNotFoundHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)

	NotFoundHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var body model.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("JSONを読めません: %v", err)
	}
	if body.Error.Code != model.CodeNotFound {
		t.Errorf("error.code = %q, want %q", body.Error.Code, model.CodeNotFound)
	}
}
