package service

import (
	"fmt"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const (
	poorCloudCover   = 70.0
	poorMaxElevation = 20.0
	goodCloudCover   = 30.0
	goodMaxElevation = 40.0
	goodMagnitude    = -2.0
)

var jst = time.FixedZone("JST", 9*60*60)

var visibilityLabels = map[string]string{
	model.VisibilityGood:    "良好",
	model.VisibilityFair:    "通常",
	model.VisibilityPoor:    "条件が厳しい",
	model.VisibilityUnknown: "判定不明",
}

func judgeVisibility(maxElevation float64, cloudCover, magnitude *float64) (string, []string) {
	var poorReasons []string
	if cloudCover != nil && *cloudCover >= poorCloudCover {
		poorReasons = append(poorReasons, fmt.Sprintf("雲量が%g%%以上", poorCloudCover))
	}
	if maxElevation < poorMaxElevation {
		poorReasons = append(poorReasons, fmt.Sprintf("最大仰角が%g度未満", poorMaxElevation))
	}
	if len(poorReasons) > 0 {
		return model.VisibilityPoor, poorReasons
	}

	var missing []string
	if cloudCover == nil {
		missing = append(missing, "雲量が不明")
	}
	if magnitude == nil {
		missing = append(missing, "明るさが不明")
	}
	if len(missing) > 0 {
		return model.VisibilityUnknown, missing
	}

	if *cloudCover < goodCloudCover && maxElevation >= goodMaxElevation && *magnitude <= goodMagnitude {
		return model.VisibilityGood, []string{
			fmt.Sprintf("最大仰角が%g度以上", goodMaxElevation),
			fmt.Sprintf("明るさが%g等級以下", goodMagnitude),
			fmt.Sprintf("雲量が%g%%未満", goodCloudCover),
		}
	}

	return model.VisibilityFair, []string{
		"良好・条件が厳しい、いずれの基準にも該当しない",
		fmt.Sprintf("雲量%g%%", *cloudCover),
		fmt.Sprintf("最大仰角%g度", maxElevation),
		fmt.Sprintf("明るさ%g等級", *magnitude),
	}
}

func visibilityMessage(start model.PassPoint, visibility string) string {
	when := fmt.Sprintf("%sごろ%sの空が観測候補です。", start.Time.In(jst).Format("15:04"), start.Direction)

	if visibility == model.VisibilityUnknown {
		return when + "天気または明るさが不明のため判定できません。"
	}
	return when + "観測おすすめ度は" + visibilityLabels[visibility] + "です。"
}
