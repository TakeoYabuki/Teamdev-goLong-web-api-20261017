package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/service"
)

func TestParsePassesQuery(t *testing.T) {
	kagoshima := model.Location{ID: "kagoshima", Name: "鹿児島市", Lat: 31.5966, Lon: 130.5571}
	tokyo := model.Location{ID: "tokyo", Name: "東京都心", Lat: 35.6812, Lon: 139.7671}
	minElevation := 30.0

	tests := []struct {
		name  string
		query string
		want  service.PassesQuery
	}{
		{
			name:  "cityだけなら既定値",
			query: "city=kagoshima",
			want:  service.PassesQuery{Location: kagoshima, SortKey: "start", Page: 1, Limit: 10},
		},
		{
			name:  "すべて指定する",
			query: "city=tokyo&visibility=good,fair&minElevation=30&sort=-maxElevation&page=2&limit=5",
			want: service.PassesQuery{
				Location:     tokyo,
				Visibilities: []string{"good", "fair"},
				MinElevation: &minElevation,
				SortKey:      "maxElevation",
				Descending:   true,
				Page:         2,
				Limit:        5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, _ := url.ParseQuery(tt.query)

			got, err := parsePassesQuery(params)
			if err != nil {
				t.Fatalf("エラーになりました: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestParsePassesQueryErrors(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"cityがない", "visibility=good"},
		{"知らない都市", "city=atlantis"},
		{"知らないパラメータ", "city=tokyo&days=10"},
		{"visibilityに知らない値", "city=tokyo&visibility=great"},
		{"visibilityに空の値", "city=tokyo&visibility=good,"},
		{"minElevationが90より大きい", "city=tokyo&minElevation=91"},
		{"minElevationが負", "city=tokyo&minElevation=-1"},
		{"minElevationが数でない", "city=tokyo&minElevation=abc"},
		{"minElevationがNaN", "city=tokyo&minElevation=NaN"},
		{"sortに知らない項目", "city=tokyo&sort=end"},
		{"sortがマイナスだけ", "city=tokyo&sort=-"},
		{"pageが0", "city=tokyo&page=0"},
		{"limitが51", "city=tokyo&limit=51"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, _ := url.ParseQuery(tt.query)

			if _, err := parsePassesQuery(params); err == nil {
				t.Error("エラーになるはずが、通ってしまいました")
			}
		})
	}
}

func TestPassesHandler(t *testing.T) {
	succeed := func(ctx context.Context, query service.PassesQuery) (model.PassesResponse, error) {
		return model.PassesResponse{Location: query.Location, Passes: []model.Pass{}}, nil
	}
	timeout := func(ctx context.Context, query service.PassesQuery) (model.PassesResponse, error) {
		return model.PassesResponse{}, context.DeadlineExceeded
	}
	broken := func(ctx context.Context, query service.PassesQuery) (model.PassesResponse, error) {
		return model.PassesResponse{}, errors.New("N2YOの失敗")
	}

	tests := []struct {
		name       string
		method     string
		query      string
		getPasses  func(context.Context, service.PassesQuery) (model.PassesResponse, error)
		wantStatus int
		wantCode   string
	}{
		{"正常", http.MethodGet, "?city=tokyo", succeed, http.StatusOK, ""},
		{"GET以外は405", http.MethodPost, "?city=tokyo", succeed, http.StatusMethodNotAllowed, model.CodeMethodNotAllowed},
		{"不正なパラメータは400", http.MethodGet, "?city=atlantis", succeed, http.StatusBadRequest, model.CodeInvalidParameter},
		{"外部のタイムアウトは504", http.MethodGet, "?city=tokyo", timeout, http.StatusGatewayTimeout, model.CodeUpstreamTimeout},
		{"外部の失敗は502", http.MethodGet, "?city=tokyo", broken, http.StatusBadGateway, model.CodeUpstreamError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/api/iss/passes"+tt.query, nil)

			passesHandler(rec, req, tt.getPasses)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantCode == "" {
				var body model.PassesResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("JSONを読めません: %v", err)
				}
				if body.Location.ID != "tokyo" {
					t.Errorf("location.id = %q, want tokyo", body.Location.ID)
				}
				return
			}
			var body model.ErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("JSONを読めません: %v", err)
			}
			if body.Error.Code != tt.wantCode {
				t.Errorf("error.code = %q, want %q", body.Error.Code, tt.wantCode)
			}
		})
	}
}
