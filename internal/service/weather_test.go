package service

import (
	"testing"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestMatchCloudCover(t *testing.T) {
	tokyo := time.FixedZone("JST", 9*60*60)
	at := func(day, hour, minute int) time.Time {
		return time.Date(2026, 10, day, hour, minute, 0, 0, tokyo)
	}
	cloud := func(v float64) *float64 {
		return &v
	}

	hours := []model.HourlyCloud{
		{Time: at(24, 11, 0), CloudCover: cloud(20)},
		{Time: at(24, 12, 0), CloudCover: cloud(40)},
		{Time: at(24, 13, 0), CloudCover: nil},
		{Time: at(24, 14, 0), CloudCover: cloud(85)},
		{Time: at(24, 23, 0), CloudCover: cloud(10)},
		{Time: at(25, 0, 0), CloudCover: cloud(60)},
	}

	tests := []struct {
		name          string
		maxTime       time.Time
		wantCloud     *float64
		wantWeatherAt time.Time
	}{
		{"ちょうどの時刻", at(24, 12, 0), cloud(40), at(24, 12, 0)},
		{"近いほうを選ぶ", at(24, 12, 20), cloud(40), at(24, 12, 0)},
		{"同じ近さなら早い時刻", at(24, 11, 30), cloud(20), at(24, 11, 0)},
		{"30分ちょうどは使う", at(24, 14, 30), cloud(85), at(24, 14, 0)},
		{"30分を超えたら不明", at(24, 14, 31), nil, time.Time{}},
		{"最初より前で30分を超えたら不明", at(24, 10, 29), nil, time.Time{}},
		{"いちばん近い予報が欠けていたら不明", at(24, 12, 40), nil, time.Time{}},
		{"日付をまたいだ予報を選ぶ", at(24, 23, 50), cloud(60), at(25, 0, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCloud, gotWeatherAt := matchCloudCover(hours, tt.maxTime)

			if tt.wantCloud == nil {
				if gotCloud != nil {
					t.Errorf("不明のはずが、雲量 %v が返りました", *gotCloud)
				}
				if gotWeatherAt != nil {
					t.Errorf("不明のはずが、予報の時刻 %v が返りました", *gotWeatherAt)
				}
				return
			}
			if gotCloud == nil || *gotCloud != *tt.wantCloud {
				t.Fatalf("雲量 = %v, want %v", gotCloud, *tt.wantCloud)
			}
			if gotWeatherAt == nil || !gotWeatherAt.Equal(tt.wantWeatherAt) {
				t.Errorf("予報の時刻 = %v, want %v", gotWeatherAt, tt.wantWeatherAt)
			}
		})
	}
}

func TestMatchCloudCoverWithoutForecast(t *testing.T) {
	gotCloud, gotWeatherAt := matchCloudCover(nil, time.Now())
	if gotCloud != nil || gotWeatherAt != nil {
		t.Error("予報がないときは不明のはずです")
	}
}

// 予報のデータが上書きされてないかの追加テスト
func TestMatchCloudCoverDoesNotModifyOriginal(t *testing.T) {
	originalCloud := 40.0
	originalTime := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)

	hours := []model.HourlyCloud{
		{
			Time:       originalTime,
			CloudCover: &originalCloud,
		},
	}

	// 予報データを取得
	gotCloud, gotWeatherAt := matchCloudCover(hours, originalTime)

	if gotCloud == nil || gotWeatherAt == nil {
		t.Fatal("雲量と予報時刻が取得できませんでした")
	}

	// 戻り値を変更する
	*gotCloud = 80.0
	*gotWeatherAt = gotWeatherAt.Add(time.Hour)

	// 元のデータが変更されていないか確認
	if *hours[0].CloudCover != 40.0 {
		t.Errorf("元の雲量が変更されました: %v", *hours[0].CloudCover)
	}

	if !hours[0].Time.Equal(originalTime) {
		t.Errorf("元の予報時刻が変更されました: %v", hours[0].Time)
	}
}
