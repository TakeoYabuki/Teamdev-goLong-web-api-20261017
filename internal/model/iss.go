package model

import "time"

// ISSの位置情報を表す構造体
type ISSPosition struct {
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	AltitudeKm  float64   `json:"altitudeKm"`
	VelocityKmh float64   `json:"velocityKmh"`
	PositionAt  time.Time `json:"positionAt"`
	FetchedAt   time.Time `json:"fetchedAt"`
}
