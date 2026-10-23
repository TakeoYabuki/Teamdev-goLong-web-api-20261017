package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/handler"
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestLocationHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		query      string
		wantStatus int
		wantCount  int
		wantPage   int
		wantLimit  int
		wantTotal  int
		wantPages  int
		wantID     string
		wantError  string
	}{
		{
			name:       "パラメータなしで最初の10件を返す",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantCount:  10,
			wantPage:   1,
			wantLimit:  10,
			wantTotal:  47,
			wantPages:  5,
		},
		{
			name:       "鹿児島を検索すると鹿児島市を返す",
			method:     http.MethodGet,
			query:      "?q=" + url.QueryEscape("鹿児島"),
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantPage:   1,
			wantLimit:  10,
			wantTotal:  1,
			wantPages:  1,
			wantID:     "kagoshima",
		},
		{
			name:       "最終ページは7件を返す",
			method:     http.MethodGet,
			query:      "?page=5&limit=10",
			wantStatus: http.StatusOK,
			wantCount:  7,
			wantPage:   5,
			wantLimit:  10,
			wantTotal:  47,
			wantPages:  5,
		},
		{
			name:       "範囲外ページは空配列を返す",
			method:     http.MethodGet,
			query:      "?page=6&limit=10",
			wantStatus: http.StatusOK,
			wantCount:  0,
			wantPage:   6,
			wantLimit:  10,
			wantTotal:  47,
			wantPages:  5,
		},
		{
			name:       "pageが0なら400を返す",
			method:     http.MethodGet,
			query:      "?page=0",
			wantStatus: http.StatusBadRequest,
			wantError:  "INVALID_PARAMETER",
		},
		{
			name:       "知らないパラメータなら400を返す",
			method:     http.MethodGet,
			query:      "?city=tokyo",
			wantStatus: http.StatusBadRequest,
			wantError:  "INVALID_PARAMETER",
		},
		{
			name:       "同じパラメータが2回なら400を返す",
			method:     http.MethodGet,
			query:      "?page=1&page=2",
			wantStatus: http.StatusBadRequest,
			wantError:  "INVALID_PARAMETER",
		},
		{
			name:       "POSTなら405を返す",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "METHOD_NOT_ALLOWED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(
				tt.method,
				"/api/locations/search"+tt.query,
				nil,
			)
			rec := httptest.NewRecorder()

			handler.LocationHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d: %s",
					rec.Code,
					tt.wantStatus,
					rec.Body.String(),
				)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf(
					"Content-Type = %q, want application/json",
					got,
				)
			}
			if tt.wantError != "" {
				var body model.ErrorResponse

				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("JSONを読めません: %v", err)
				}
				if body.Error.Code != tt.wantError {
					t.Errorf(
						"error.code = %q, want %q",
						body.Error.Code,
						tt.wantError,
					)
				}
				if body.Error.Message == "" {
					t.Error("エラーメッセージが空です")
				}
				if tt.wantStatus == http.StatusMethodNotAllowed {
					if got := rec.Header().Get("Allow"); got != http.MethodGet {
						t.Errorf("Allow = %q, want GET", got)
					}
				}
				return
			}
			// 成功の場合は都市検索レスポンスを確認する
			var body model.LocationSearchResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("JSONを読めません: %v", err)
			}
			if body.Results == nil {
				t.Error("resultsがnullです。配列を返してください")
			}
			if len(body.Results) != tt.wantCount {
				t.Fatalf(
					"件数 = %d, want %d",
					len(body.Results),
					tt.wantCount,
				)
			}
			if body.Page != tt.wantPage ||
				body.Limit != tt.wantLimit ||
				body.Total != tt.wantTotal ||
				body.TotalPages != tt.wantPages {

				t.Errorf(
					"ページ情報 = page:%d limit:%d total:%d totalPages:%d, "+
						"want page:%d limit:%d total:%d totalPages:%d",
					body.Page, body.Limit, body.Total, body.TotalPages,
					tt.wantPage, tt.wantLimit, tt.wantTotal, tt.wantPages,
				)
			}
			// 検索のケースでは帰った都市のIDも確認する
			if tt.wantID != "" {
				if len(body.Results) == 0 {
					t.Fatal("都市が返っていません")
				}
				if body.Results[0].ID != tt.wantID {
					t.Errorf(
						"都市ID = %q, want %q",
						body.Results[0].ID,
						tt.wantID,
					)
				}
			}
		})
	}
}
