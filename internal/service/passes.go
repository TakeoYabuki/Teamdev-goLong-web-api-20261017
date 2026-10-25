package service

import "github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"

const (
	SortByStart        = "start"
	SortByMaxElevation = "maxElevation"
	SortByCloudCover   = "cloudCover"
)

type PassesQuery struct {
	Location     model.Location
	Visibilities []string
	MinElevation *float64
	SortKey      string
	Descending   bool
	Page         int
	Limit        int
}
