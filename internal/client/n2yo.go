package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const n2yoBaseURL = "https://api.n2yo.com/rest/v1/satellite/visualpasses"

func FetchVisualPasses(ctx context.Context, lat, lon float64) (*model.N2YOResponse, error) {
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("N2YO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("N2YO_API_KEYが設定されていません")
	}

	// 緯度・経度を文字列に変換
	latitude := strconv.FormatFloat(lat, 'f', -1, 64)
	longitude := strconv.FormatFloat(lon, 'f', -1, 64)

	// N2YOのURLを作成
	requestURL := fmt.Sprintf(
		"%s/25544/%s/%s/0/4/1/",
		n2yoBaseURL,
		latitude,
		longitude,
	)

	// APIキーをクエリパラメータとして追加
	query := url.Values{}
	query.Set("apiKey", apiKey)
	requestURL += "&" + query.Encode()

	// HTTPリクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	// HTTPクライアントを作成
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// N2YOへリクエストを送信
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("N2YOが異常なステータスを返しました: %d", resp.StatusCode)
	}

	var raw json.RawMessage

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	// N2YOがHTTP 200でエラーJSONを返した場合の確認
	var apiError struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(raw, &apiError); err != nil {
		return nil, err
	}

	if apiError.Error != "" {
		return nil, fmt.Errorf("N2YO APIエラー: %s", apiError.Error)
	}

	// 正常なJSONをN2YOResponseに変換

	var data model.N2YOResponse

	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
