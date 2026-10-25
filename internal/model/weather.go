package model

import "time"

type CloudForecast struct {
	Hours     []HourlyCloud
	FetchedAt time.Time
}

type HourlyCloud struct {
	Time       time.Time
	CloudCover *float64
}
