package service

import (
	"slices"
	"testing"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func ptr(v float64) *float64 {
	return &v
}

func TestJudgeVisibility(t *testing.T) {
	tests := []struct {
		name         string
		maxElevation float64
		cloudCover   *float64
		magnitude    *float64
		want         string
	}{
		{"すべて良ければgood", 62, ptr(20), ptr(-3.1), model.VisibilityGood},
		{"雲量70ちょうどはpoor", 62, ptr(70), ptr(-3), model.VisibilityPoor},
		{"雲量69.9はpoorではない", 62, ptr(69.9), ptr(-3), model.VisibilityFair},
		{"仰角20ちょうどはpoorではない", 20, ptr(10), ptr(-3), model.VisibilityFair},
		{"仰角19.9はpoor", 19.9, ptr(10), ptr(-3), model.VisibilityPoor},
		{"雲量30ちょうどはgoodではない", 62, ptr(30), ptr(-3), model.VisibilityFair},
		{"雲量29.9はgood", 62, ptr(29.9), ptr(-3), model.VisibilityGood},
		{"仰角40ちょうどはgood", 40, ptr(10), ptr(-3), model.VisibilityGood},
		{"仰角39.9はgoodではない", 39.9, ptr(10), ptr(-3), model.VisibilityFair},
		{"明るさ-2.0ちょうどはgood", 62, ptr(10), ptr(-2.0), model.VisibilityGood},
		{"明るさ-1.9はgoodではない", 62, ptr(10), ptr(-1.9), model.VisibilityFair},
		{"雲量が不明ならunknown", 62, nil, ptr(-3), model.VisibilityUnknown},
		{"明るさが不明ならunknown", 62, ptr(10), nil, model.VisibilityUnknown},
		{"雲量が不明でも仰角が低ければpoor", 10, nil, ptr(-3), model.VisibilityPoor},
		{"明るさが不明でも雲量が多ければpoor", 62, ptr(80), nil, model.VisibilityPoor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := judgeVisibility(tt.maxElevation, tt.cloudCover, tt.magnitude)
			if got != tt.want {
				t.Errorf("visibility = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJudgeVisibilityReasons(t *testing.T) {
	tests := []struct {
		name         string
		maxElevation float64
		cloudCover   *float64
		magnitude    *float64
		want         []string
	}{
		{"goodは3条件", 62, ptr(20), ptr(-3.1), []string{"最大仰角が40度以上", "明るさが-2等級以下", "雲量が30%未満"}},
		{"poorは当てはまった悪条件", 10, ptr(80), ptr(-3), []string{"雲量が70%以上", "最大仰角が20度未満"}},
		{"unknownは足りない項目", 62, nil, nil, []string{"雲量が不明", "明るさが不明"}},
		{"fairは取れた値", 35, ptr(45), ptr(-1.5), []string{"良好・条件が厳しい、いずれの基準にも該当しない", "雲量45%", "最大仰角35度", "明るさ-1.5等級"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := judgeVisibility(tt.maxElevation, tt.cloudCover, tt.magnitude)
			if !slices.Equal(got, tt.want) {
				t.Errorf("reasons = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVisibilityMessage(t *testing.T) {
	start := model.PassPoint{
		Time:      time.Date(2026, 10, 19, 9, 42, 10, 0, time.UTC),
		Direction: "北西",
	}

	tests := []struct {
		visibility string
		want       string
	}{
		{model.VisibilityGood, "18:42ごろ北西の空が観測候補です。観測おすすめ度は良好です。"},
		{model.VisibilityUnknown, "18:42ごろ北西の空が観測候補です。天気または明るさが不明のため判定できません。"},
	}
	for _, tt := range tests {
		t.Run(tt.visibility, func(t *testing.T) {
			if got := visibilityMessage(start, tt.visibility); got != tt.want {
				t.Errorf("message = %q, want %q", got, tt.want)
			}
		})
	}
}
