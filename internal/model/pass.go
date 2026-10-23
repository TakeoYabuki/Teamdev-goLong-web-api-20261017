package model

import "time"

type PassesResponse struct {
	Location             Location    `json:"location"`
	From                 time.Time   `json:"from"`
	To                   time.Time   `json:"to"`
	NextCandidate        *Pass       `json:"nextCandidate"`
	Summary              PassSummary `json:"summary"`
	Passes               []Pass      `json:"passes"`
	Page                 int         `json:"page"`
	Limit                int         `json:"limit"`
	Total                int         `json:"total"`
	TotalPages           int         `json:"totalPages"`
	PredictionsFetchedAt time.Time   `json:"predictionsFetchedAt"`
	WeatherFetchedAt     *time.Time  `json:"weatherFetchedAt"`
}

type PassSummary struct {
	Total   int `json:"total"`
	Good    int `json:"good"`
	Fair    int `json:"fair"`
	Poor    int `json:"poor"`
	Unknown int `json:"unknown"`
}

type Pass struct {
	ID          string     `json:"id"`
	Start       PassPoint  `json:"start"`
	Max         PassPoint  `json:"max"`
	End         PassPoint  `json:"end"`
	DurationSec int        `json:"durationSec"`
	Magnitude   *float64   `json:"magnitude"`
	CloudCover  *float64   `json:"cloudCover"`
	WeatherAt   *time.Time `json:"weatherAt"`
	Visibility  string     `json:"visibility"`
	Reasons     []string   `json:"reasons"`
	Message     string     `json:"message"`
}

type PassPoint struct {
	Time      time.Time `json:"time"`
	Direction string    `json:"direction"`
	Azimuth   float64   `json:"azimuth"`
	Elevation float64   `json:"elevation"`
}

const (
	VisibilityGood    = "good"
	VisibilityFair    = "fair"
	VisibilityPoor    = "poor"
	VisibilityUnknown = "unknown"
)
