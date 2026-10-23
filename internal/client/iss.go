package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

// 外部APIのJSONを受け取るための構造体
type issAPIResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Velocity  float64 `json:"velocity"`
	Timestamp int64   `json:"timestamp"`
}

// 本番用の入口
func FetchISSPosition(ctx context.Context) (*model.ISSPosition, error) {
	return fetchISSPosition(
		ctx,
		"https://api.wheretheiss.at/v1/satellites/25544",
	)
}

// 実際にHTTP通信する関数
func fetchISSPosition(ctx context.Context, url string) (*model.ISSPosition, error) {

	// GETリクエストを作成
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}
	// HTTPクライアントを作成 。タイムアウトを１０秒に設定
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 外部APIへ送信
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// HTTPステータスを確認
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ISS外部APIが異常なステータスを返しました: %d", resp.StatusCode)
	}

	// JSONをデコード
	var data issAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// modelの構造体に変換
	position := model.ISSPosition{
		Lat:         data.Latitude,
		Lon:         data.Longitude,
		AltitudeKm:  data.Altitude,
		VelocityKmh: data.Velocity,
		PositionAt:  time.Unix(data.Timestamp, 0),
		FetchedAt:   time.Now().Truncate(time.Second),
	}

	return &position, nil
}
