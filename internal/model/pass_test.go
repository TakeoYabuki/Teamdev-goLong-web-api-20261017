package model

import (
	"encoding/json"
	"testing"
	"time"
)

var jst = time.FixedZone("JST", 9*60*60)

func TestPassJSONMatchesSpec(t *testing.T) {
	magnitude := -3.1
	cloudCover := 20.0
	weatherAt := time.Date(2026, 10, 19, 19, 0, 0, 0, jst)

	pass := Pass{
		ID: "kagoshima_1792402930",
		Start: PassPoint{
			Time:      time.Date(2026, 10, 19, 18, 42, 10, 0, jst),
			Direction: "北西",
			Azimuth:   315,
			Elevation: 10,
		},
		Max: PassPoint{
			Time:      time.Date(2026, 10, 19, 18, 45, 20, 0, jst),
			Direction: "北東",
			Azimuth:   45,
			Elevation: 62,
		},
		End: PassPoint{
			Time:      time.Date(2026, 10, 19, 18, 48, 35, 0, jst),
			Direction: "南東",
			Azimuth:   135,
			Elevation: 10,
		},
		DurationSec: 385,
		Magnitude:   &magnitude,
		CloudCover:  &cloudCover,
		WeatherAt:   &weatherAt,
		Visibility:  VisibilityGood,
		Reasons:     []string{"最大仰角が40度以上", "明るさが-2等級以下", "雲量が30%未満"},
		Message:     "18:42ごろ北西の空が観測候補です。観測おすすめ度は良好です。",
	}

	want := `{"id":"kagoshima_1792402930",` +
		`"start":{"time":"2026-10-19T18:42:10+09:00","direction":"北西","azimuth":315,"elevation":10},` +
		`"max":{"time":"2026-10-19T18:45:20+09:00","direction":"北東","azimuth":45,"elevation":62},` +
		`"end":{"time":"2026-10-19T18:48:35+09:00","direction":"南東","azimuth":135,"elevation":10},` +
		`"durationSec":385,"magnitude":-3.1,"cloudCover":20,"weatherAt":"2026-10-19T19:00:00+09:00",` +
		`"visibility":"good","reasons":["最大仰角が40度以上","明るさが-2等級以下","雲量が30%未満"],` +
		`"message":"18:42ごろ北西の空が観測候補です。観測おすすめ度は良好です。"}`

	got, err := json.Marshal(pass)
	if err != nil {
		t.Fatalf("JSONにできません: %v", err)
	}
	if string(got) != want {
		t.Errorf("JSONが仕様書5章と違います\n got: %s\nwant: %s", got, want)
	}
}

func TestPassesResponseJSONMatchesSpec(t *testing.T) {
	response := PassesResponse{
		Location: Location{ID: "kagoshima", Name: "鹿児島市", Lat: 31.5966, Lon: 130.5571},
		From:     time.Date(2026, 10, 19, 15, 0, 0, 0, jst),
		To:       time.Date(2026, 10, 22, 15, 0, 0, 0, jst),
		Summary:  PassSummary{},
		Passes:   []Pass{},
		Page:     1,
		Limit:    5,

		PredictionsFetchedAt: time.Date(2026, 10, 19, 14, 45, 0, 0, jst),
	}

	want := `{"location":{"id":"kagoshima","name":"鹿児島市","lat":31.5966,"lon":130.5571},` +
		`"from":"2026-10-19T15:00:00+09:00","to":"2026-10-22T15:00:00+09:00",` +
		`"nextCandidate":null,` +
		`"summary":{"total":0,"good":0,"fair":0,"poor":0,"unknown":0},` +
		`"passes":[],"page":1,"limit":5,"total":0,"totalPages":0,` +
		`"predictionsFetchedAt":"2026-10-19T14:45:00+09:00","weatherFetchedAt":null}`

	got, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("JSONにできません: %v", err)
	}
	if string(got) != want {
		t.Errorf("JSONが仕様書6章と違います\n got: %s\nwant: %s", got, want)
	}
}
