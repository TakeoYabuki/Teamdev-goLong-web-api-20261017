package service

import (
	"cmp"
	"slices"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const (
	SortByStart        = "start"
	SortByMaxElevation = "maxElevation"
	SortByCloudCover   = "cloudCover"
)

const passesWindow = 72 * time.Hour

type PassesQuery struct {
	Location     model.Location
	Visibilities []string
	MinElevation *float64
	SortKey      string
	Descending   bool
	Page         int
	Limit        int
}

func buildPassesResponse(
	query PassesQuery,
	passes []model.Pass,
	predictionsFetchedAt time.Time,
	forecast *model.CloudForecast,
	now time.Time,
) model.PassesResponse {
	from := now.In(jst).Truncate(time.Second)
	to := from.Add(passesWindow)

	var hours []model.HourlyCloud
	var weatherFetchedAt *time.Time
	if forecast != nil {
		hours = forecast.Hours
		fetchedAt := forecast.FetchedAt.In(jst)
		weatherFetchedAt = &fetchedAt
	}

	judged := judgePasses(passes, hours, from, to)

	var nextCandidate *model.Pass
	for i := range judged {
		if nextCandidate == nil || judged[i].Start.Time.Before(nextCandidate.Start.Time) {
			nextCandidate = &judged[i]
		}
	}

	filtered := filterPasses(judged, query)
	sortPasses(filtered, query.SortKey, query.Descending)

	total := len(filtered)
	pagePasses := []model.Pass{}
	if start := (query.Page - 1) * query.Limit; start < total {
		pagePasses = filtered[start:min(start+query.Limit, total)]
	}

	return model.PassesResponse{
		Location:             query.Location,
		From:                 from,
		To:                   to,
		NextCandidate:        nextCandidate,
		Summary:              summarizePasses(judged),
		Passes:               pagePasses,
		Page:                 query.Page,
		Limit:                query.Limit,
		Total:                total,
		TotalPages:           (total + query.Limit - 1) / query.Limit,
		PredictionsFetchedAt: predictionsFetchedAt.In(jst).Truncate(time.Second),
		WeatherFetchedAt:     weatherFetchedAt,
	}
}

func judgePasses(passes []model.Pass, hours []model.HourlyCloud, from, to time.Time) []model.Pass {
	judged := []model.Pass{}

	for _, pass := range passes {
		if pass.Start.Time.Before(from) || !pass.Start.Time.Before(to) {
			continue
		}

		cloudCover, weatherAt := matchCloudCover(hours, pass.Max.Time)
		if weatherAt != nil {
			inJST := weatherAt.In(jst)
			weatherAt = &inJST
		}

		pass.CloudCover = cloudCover
		pass.WeatherAt = weatherAt
		pass.Visibility, pass.Reasons = judgeVisibility(pass.Max.Elevation, cloudCover, pass.Magnitude)
		pass.Message = visibilityMessage(pass.Start, pass.Visibility)

		judged = append(judged, pass)
	}
	return judged
}

func summarizePasses(passes []model.Pass) model.PassSummary {
	summary := model.PassSummary{Total: len(passes)}

	for _, pass := range passes {
		switch pass.Visibility {
		case model.VisibilityGood:
			summary.Good++
		case model.VisibilityFair:
			summary.Fair++
		case model.VisibilityPoor:
			summary.Poor++
		case model.VisibilityUnknown:
			summary.Unknown++
		}
	}
	return summary
}

func filterPasses(passes []model.Pass, query PassesQuery) []model.Pass {
	filtered := []model.Pass{}

	for _, pass := range passes {
		if len(query.Visibilities) > 0 && !slices.Contains(query.Visibilities, pass.Visibility) {
			continue
		}
		if query.MinElevation != nil && pass.Max.Elevation < *query.MinElevation {
			continue
		}
		filtered = append(filtered, pass)
	}
	return filtered
}

func sortPasses(passes []model.Pass, key string, descending bool) {
	slices.SortStableFunc(passes, func(a, b model.Pass) int {
		if key == SortByCloudCover && (a.CloudCover == nil) != (b.CloudCover == nil) {
			if a.CloudCover == nil {
				return 1
			}
			return -1
		}

		var order int
		switch key {
		case SortByMaxElevation:
			order = cmp.Compare(a.Max.Elevation, b.Max.Elevation)
		case SortByCloudCover:
			if a.CloudCover != nil && b.CloudCover != nil {
				order = cmp.Compare(*a.CloudCover, *b.CloudCover)
			}
		default:
			order = a.Start.Time.Compare(b.Start.Time)
		}
		if descending {
			order = -order
		}
		if order != 0 {
			return order
		}

		if order := a.Start.Time.Compare(b.Start.Time); order != 0 {
			return order
		}
		return cmp.Compare(a.ID, b.ID)
	})
}
