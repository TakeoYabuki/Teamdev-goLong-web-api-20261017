package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestFetchCloudForecast(t *testing.T) {
	body, err := os.ReadFile("testdata/open_meteo_forecast.json")
	if err != nil {
		t.Fatalf("テスト用の応答を読めません: %v", err)
	}

	var query url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Write(body)
	}))
	defer server.Close()

	forecast, err := fetchCloudForecast(context.Background(), server.URL, 35.6812, 139.7671)
	if err != nil {
		t.Fatalf("取得に失敗しました: %v", err)
	}

	wantQuery := map[string]string{
		"latitude":       "35.6812",
		"longitude":      "139.7671",
		"hourly":         "cloud_cover",
		"forecast_hours": "76",
		"timeformat":     "unixtime",
	}
	for name, want := range wantQuery {
		if got := query.Get(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}

	if len(forecast.Hours) != 4 {
		t.Fatalf("時間の数 = %d, want 4", len(forecast.Hours))
	}
	if want := time.Unix(1792807200, 0); !forecast.Hours[0].Time.Equal(want) {
		t.Errorf("1つ目の時刻 = %v, want %v", forecast.Hours[0].Time, want)
	}
	if got := forecast.Hours[0].CloudCover; got == nil || *got != 20 {
		t.Errorf("1つ目の雲量 = %v, want 20", got)
	}
	if got := forecast.Hours[1].CloudCover; got != nil {
		t.Errorf("nullの雲量 = %v, want nil", *got)
	}
	if got := forecast.Hours[2].CloudCover; got != nil {
		t.Errorf("範囲外(120)の雲量 = %v, want nil", *got)
	}
	if forecast.FetchedAt.IsZero() {
		t.Error("取得時刻が入っていません")
	}
}

func TestFetchCloudForecastErrors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{"サーバの異常", http.StatusInternalServerError, ""},
		{"JSONが壊れている", http.StatusOK, `{"hourly":`},
		{"時刻と値の数が合わない", http.StatusOK, `{"hourly":{"time":[1792807200,1792810800],"cloud_cover":[20]}}`},
		{"データが空", http.StatusOK, `{"hourly":{"time":[],"cloud_cover":[]}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			forecast, err := fetchCloudForecast(context.Background(), server.URL, 35.6812, 139.7671)
			if err == nil {
				t.Fatal("エラーになるはずが、成功しました")
			}
			if forecast != nil {
				t.Errorf("エラーのときは nil のはずが、%v が返りました", forecast)
			}
		})
	}
}
