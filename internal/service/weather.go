package service

import (
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const maxWeatherGap = 30 * time.Minute

func matchCloudCover(hours []model.HourlyCloud, at time.Time) (*float64, *time.Time) {
	var nearest *model.HourlyCloud
	var nearestGap time.Duration

	for i := range hours {
		gap := hours[i].Time.Sub(at).Abs()
		if nearest == nil || gap < nearestGap {
			nearest = &hours[i]
			nearestGap = gap
		}
	}

	if nearest == nil || nearestGap > maxWeatherGap || nearest.CloudCover == nil {
		return nil, nil
	}

	cloudCover := *nearest.CloudCover
	weatherAt := nearest.Time
	return &cloudCover, &weatherAt
}
