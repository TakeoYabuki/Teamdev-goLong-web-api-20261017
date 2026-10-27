package service

import (
	"testing"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestGetVisualPassesConversion(t *testing.T) {
	// N2YOから取得したと仮定するテストデータ
	data := &model.N2YOResponse{
		Passes: []model.N2YOPass{
			{
				StartUTC:       1792523360,
				StartAz:        274.66,
				StartAzCompass: "W",
				StartEl:        12.02,

				MaxUTC:       1792523635,
				MaxAz:        333.04,
				MaxAzCompass: "NNW",
				MaxEl:        13.72,

				EndUTC:       1792523920,
				EndAz:        34.33,
				EndAzCompass: "NE",
				EndEl:        0.27,

				Duration: 210,
				Mag:      ptr(0.2),
			},
		},
	}

	// serviceの変換関数を呼び出す
	passes, err := getVisualPasses("kagoshima", data)
	if err != nil {
		t.Fatalf("変換エラー: %v", err)
	}

	// 通過件数
	if len(passes) != 1 {
		t.Fatalf("通過件数: got %d, want 1", len(passes))
	}

	p := passes[0]

	// ID
	if p.ID != "kagoshima_1792523360" {
		t.Errorf("ID: got %s", p.ID)
	}

	// 開始・最高点・終了時刻
	jst := time.FixedZone("JST", 9*60*60)

	if !p.Start.Time.Equal(time.Unix(1792523360, 0)) {
		t.Errorf("Start.Timeが一致しません")
	}
	if !p.Max.Time.Equal(time.Unix(1792523635, 0)) {
		t.Errorf("Max.Timeが一致しません")
	}
	if !p.End.Time.Equal(time.Unix(1792523920, 0)) {
		t.Errorf("End.Timeが一致しません")
	}

	// 日本時間（+09:00）になっているか
	if p.Start.Time.Format(time.RFC3339) != time.Unix(1792523360, 0).In(jst).Format(time.RFC3339) {
		t.Errorf("Start.Timeが日本時間になっていません")
	}

	// 方位変換
	if p.Start.Direction != "西" {
		t.Errorf("Start.Direction: got %s, want 西", p.Start.Direction)
	}
	if p.Max.Direction != "北西" {
		t.Errorf("Max.Direction: got %s, want 北西", p.Max.Direction)
	}
	if p.End.Direction != "北東" {
		t.Errorf("End.Direction: got %s, want 北東", p.End.Direction)
	}

	// 方位角と仰角
	if p.Start.Azimuth != 274.66 || p.Start.Elevation != 12.02 {
		t.Errorf("Startの方位角または仰角が一致しません")
	}
	if p.Max.Azimuth != 333.04 || p.Max.Elevation != 13.72 {
		t.Errorf("Maxの方位角または仰角が一致しません")
	}
	if p.End.Azimuth != 34.33 || p.End.Elevation != 0.27 {
		t.Errorf("Endの方位角または仰角が一致しません")
	}

	// 可視時間
	if p.DurationSec != 210 {
		t.Errorf("DurationSec: got %d, want 210", p.DurationSec)
	}

	// 明るさ
	if p.Magnitude == nil {
		t.Fatal("Magnitude: got nil, want 0.2")
	}

	if *p.Magnitude != 0.2 {
		t.Errorf("Magnitude: got %v, want 0.2", *p.Magnitude)
	}
}

// N2YOが通過0件を返した場合
func TestGetVisualPassesEmpty(t *testing.T) {
	data := &model.N2YOResponse{
		Passes: []model.N2YOPass{},
	}

	passes, err := getVisualPasses("kagoshima", data)
	if err != nil {
		t.Fatalf("変換エラー: %v", err)
	}

	if len(passes) != 0 {
		t.Errorf("通過件数: got %d, want 0", len(passes))
	}
}

// 明るさが不明値の場合
func TestGetVisualPassesUnknownMagnitude(t *testing.T) {
	unknownMag := 100000.0

	data := &model.N2YOResponse{
		Passes: []model.N2YOPass{
			{
				StartUTC: 1792523360,
				Mag:      &unknownMag,
			},
		},
	}

	passes, err := getVisualPasses("kagoshima", data)
	if err != nil {
		t.Fatal(err)
	}

	if len(passes) != 1 {
		t.Fatalf("passes: got %d, want 1", len(passes))
	}

	if passes[0].Magnitude != nil {
		t.Errorf("Magnitude: got %v, want nil", *passes[0].Magnitude)
	}
}

func TestGetVisualPassesNilMagnitude(t *testing.T) {
	data := &model.N2YOResponse{
		Passes: []model.N2YOPass{
			{
				StartUTC: 1792523360,
				Mag:      nil,
			},
		},
	}

	passes, err := getVisualPasses("kagoshima", data)
	if err != nil {
		t.Fatal(err)
	}

	if len(passes) != 1 {
		t.Fatalf("passes: got %d, want 1", len(passes))
	}

	if passes[0].Magnitude != nil {
		t.Errorf("Magnitude: got %v, want nil", *passes[0].Magnitude)
	}
}
