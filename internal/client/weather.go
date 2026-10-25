package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const (
	openMeteoURL  = "https://api.open-meteo.com/v1/forecast"
	forecastHours = 76
)

type openMeteoResponse struct {
	Hourly struct {
		Time       []int64    `json:"time"`
		CloudCover []*float64 `json:"cloud_cover"`
	} `json:"hourly"`
}

func FetchCloudForecast(ctx context.Context, lat, lon float64) (*model.CloudForecast, error) {
	return fetchCloudForecast(ctx, openMeteoURL, lat, lon)
}

func fetchCloudForecast(ctx context.Context, baseURL string, lat, lon float64) (*model.CloudForecast, error) {
	query := url.Values{}
	query.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	query.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	query.Set("hourly", "cloud_cover")
	query.Set("forecast_hours", strconv.Itoa(forecastHours))
	query.Set("timeformat", "unixtime")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Open-Meteoが異常なステータスを返しました: %d", resp.StatusCode)
	}

	var data openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	times := data.Hourly.Time
	clouds := data.Hourly.CloudCover
	if len(times) == 0 || len(times) != len(clouds) {
		return nil, errors.New("Open-Meteoの雲量の時刻と値の数が合いません")
	}

	hours := make([]model.HourlyCloud, 0, len(times))
	for i, unix := range times {
		cloud := clouds[i]
		if cloud != nil && (*cloud < 0 || *cloud > 100) {
			cloud = nil
		}
		hours = append(hours, model.HourlyCloud{Time: time.Unix(unix, 0), CloudCover: cloud})
	}

	return &model.CloudForecast{
		Hours:     hours,
		FetchedAt: time.Now().Truncate(time.Second),
	}, nil
}
