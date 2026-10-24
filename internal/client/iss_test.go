package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 外部ISS APIの代わりになるテスト用サーバー
func TestFetchISSPosition(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(`{
                "latitude": 35.68,
                "longitude": 139.69,
                "altitude": 420.5,
                "velocity": 27580.0,
                "timestamp": 1792390750
            }`))
		},
	))
	defer server.Close()

	// テスト用サーバーからISS位置情報を取得
	position, err := fetchISSPosition(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("位置情報の取得に失敗: %v", err)
	}

	// 取得した値が期待どおりか確認
	if position.Lat != 35.68 {
		t.Errorf("緯度: 期待値35.68、実際%f", position.Lat)
	}

	if position.Lon != 139.69 {
		t.Errorf("経度: 期待値139.69、実際%f", position.Lon)
	}

	if position.AltitudeKm != 420.5 {
		t.Errorf("高度: 期待値420.5、実際%f", position.AltitudeKm)
	}

	if position.VelocityKmh != 27580.0 {
		t.Errorf("速度: 期待値27580.0、実際%f", position.VelocityKmh)
	}
}

// HTTP 500を返すテスト用サーバー
func TestFetchISSPosition_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	// テスト用サーバーにアクセス
	position, err := fetchISSPosition(
		context.Background(),
		server.URL,
	)

	// エラーが返ることを確認
	if err == nil {
		t.Fatal("HTTP 500なのにエラーが返されませんでした")
	}

	// 位置情報が返されないことを確認
	if position != nil {
		t.Errorf("エラー時はnilを期待しましたが、実際は%v", position)
	}
}

// 不正なJSONを返すテスト用サーバー
func TestFetchISSPosition_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"latitude": 35.68,`))
		},
	))
	defer server.Close()

	position, err := fetchISSPosition(
		context.Background(),
		server.URL,
	)

	// JSONの解析エラーが発生することを確認
	if err == nil {
		t.Fatal("不正なJSONなのにエラーが返されませんでした")
	}

	// エラー時は位置情報が返らないことを確認
	if position != nil {
		t.Errorf("nilを期待しましたが、実際は%v", position)
	}
}

// 外部APIの代わりに空のJSONを返すサーバー
func TestFetchISSPosition_MissingFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{}`)
		},
	))
	defer server.Close()

	// テスト用サーバーからISS位置情報を取得
	position, err := fetchISSPosition(context.Background(), server.URL)

	// 必須項目が欠損しているため、エラーになることを確認
	if err == nil {
		t.Fatal("必須項目が欠損しているのにエラーが返されませんでした")
	}

	// 正常な位置情報が返されていないことを確認
	if position != nil {
		t.Errorf("positionはnilを期待しましたが、実際: %+v", position)
	}
}
