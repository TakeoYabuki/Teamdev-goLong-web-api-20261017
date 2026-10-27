package service

import (
	"context"
	"fmt"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/client"
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

// 方位角を日本語8方位に変換
func convertDirection(azimuth float64) string {
	directions := []string{
		"北", "北東", "東", "南東",
		"南", "南西", "西", "北西",
	}

	index := int((azimuth+22.5)/45) % 8

	return directions[index]
}

func GetVisualPasses(ctx context.Context, cityID string, lat, lon float64) ([]model.Pass, error) {

	//clientからN2YOのデータを取得
	data, err := client.FetchVisualPasses(ctx, lat, lon)
	if err != nil {
		return nil, err
	}

	return getVisualPasses(cityID, data)
}

func getVisualPasses(cityID string, data *model.N2YOResponse) ([]model.Pass, error) {

	//日本時間を設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	//変換後のデータを格納するスライス
	passes := make([]model.Pass, 0, len(data.Passes))

	// N2YOの通過予測を1件ずつ変換
	for _, p := range data.Passes {

		// 明るさ
		var magnitude *float64

		if p.Mag != nil && *p.Mag != 100000 {
			magnitude = p.Mag
		}

		// 共通形式のPassを作成
		pass := model.Pass{
			ID: fmt.Sprintf("%s_%d", cityID, p.StartUTC),

			Start: model.PassPoint{
				Time:      time.Unix(p.StartUTC, 0).In(jst),
				Direction: convertDirection(p.StartAz),
				Azimuth:   p.StartAz,
				Elevation: p.StartEl,
			},

			Max: model.PassPoint{
				Time:      time.Unix(p.MaxUTC, 0).In(jst),
				Direction: convertDirection(p.MaxAz),
				Azimuth:   p.MaxAz,
				Elevation: p.MaxEl,
			},

			End: model.PassPoint{
				Time:      time.Unix(p.EndUTC, 0).In(jst),
				Direction: convertDirection(p.EndAz),
				Azimuth:   p.EndAz,
				Elevation: p.EndEl,
			},

			DurationSec: p.Duration,
			Magnitude:   magnitude,
		}

		//スライスに追加
		passes = append(passes, pass)
	}

	return passes, nil
}
