package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestIssNowHandler_Post(t *testing.T) {
	// テスト用のPOSTリクエストを作成
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/iss/now",
		nil,
	)

	// テスト用のレスポンス記録器を作成
	rec := httptest.NewRecorder()

	// テスト対象のhandlerを実行
	IssNowHandler(rec, req)

	// HTTPステータスコードを確認
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"期待するステータス: 405、実際: %d",
			rec.Code,
		)
	}

	// Allowヘッダーを確認
	if rec.Header().Get("Allow") != http.MethodGet {
		t.Errorf(
			"AllowヘッダーがGETではありません: %s",
			rec.Header().Get("Allow"),
		)
	}

	// レスポンスのJSONを確認
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("JSONの解析に失敗: %v", err)
	}

	if body.Error.Code != "METHOD_NOT_ALLOWED" {
		t.Errorf("期待するエラーコード: METHOD_NOT_ALLOWED、実際: %s", body.Error.Code)
	}

	if body.Error.Message != "GETを使用してください" {
		t.Errorf("期待するメッセージ: GETを使用してください、実際: %s", body.Error.Message)
	}
}

func TestIssNowHandler_Get(t *testing.T) {
	//テスト用の関数を定義
	fakeGetPosition := func(ctx context.Context) (*model.ISSPosition, error) {
		return &model.ISSPosition{
			Lat:         35.0,
			Lon:         139.0,
			AltitudeKm:  420.0,
			VelocityKmh: 27500.0,
		}, nil
	}

	// GETリクエストを作成
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/iss/now",
		nil,
	)

	// レスポンスを記録
	rec := httptest.NewRecorder()

	// テスト用の関数を渡す
	issNowHandler(rec, req, fakeGetPosition)

	// HTTP 200が返るか確認
	if rec.Code != http.StatusOK {
		t.Fatalf(
			"期待するステータス: 200、実際: %d、レスポンス: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	// JSONの形式を確認
	var position struct {
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		AltitudeKm  float64 `json:"altitudeKm"`
		VelocityKmh float64 `json:"velocityKmh"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &position); err != nil {
		t.Fatalf("JSONの解析に失敗: %v", err)
	}

	// 緯度・経度が有効な範囲か確認
	if position.Lat < -90 || position.Lat > 90 {
		t.Errorf("緯度が範囲外です: %f", position.Lat)
	}

	if position.Lon < -180 || position.Lon > 180 {
		t.Errorf("経度が範囲外です: %f", position.Lon)
	}

	if position.AltitudeKm <= 0 {
		t.Errorf("高度が不正です: %f", position.AltitudeKm)
	}

	if position.VelocityKmh <= 0 {
		t.Errorf("速度が不正です: %f", position.VelocityKmh)
	}
}

// 外部API取得エラー発生時502を返すか
func TestIssNowHandler_UpstreamError(t *testing.T) {
	fakeGetPosition := func(ctx context.Context) (*model.ISSPosition, error) {
		return nil, errors.New("外部APIのレスポンスが不正です")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/iss/now", nil)
	rec := httptest.NewRecorder()

	issNowHandler(rec, req, fakeGetPosition)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("期待するステータス: 502、実際: %d", rec.Code)
	}
}
