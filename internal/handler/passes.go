package handler

import (
	"context"
	"errors"
	"log"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/service"
)

var visibilityValues = []string{
	model.VisibilityGood,
	model.VisibilityFair,
	model.VisibilityPoor,
	model.VisibilityUnknown,
}

var sortKeys = []string{
	service.SortByStart,
	service.SortByMaxElevation,
	service.SortByCloudCover,
}

func PassesHandler(w http.ResponseWriter, r *http.Request) {
	passesHandler(w, r, service.GetPasses)
}

func passesHandler(
	w http.ResponseWriter,
	r *http.Request,
	getPasses func(context.Context, service.PassesQuery) (model.PassesResponse, error),
) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		WriteError(w, http.StatusMethodNotAllowed, model.CodeMethodNotAllowed, "GETを使用してください")
		return
	}

	query, err := parsePassesQuery(r.URL.Query())
	if err != nil {
		WriteError(w, http.StatusBadRequest, model.CodeInvalidParameter, err.Error())
		return
	}

	response, err := getPasses(r.Context(), query)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if isUpstreamTimeout(err) {
			WriteError(w, http.StatusGatewayTimeout, model.CodeUpstreamTimeout, "通過予測の取得がタイムアウトしました")
			return
		}
		log.Printf("通過予測の取得エラー: %v", err)
		WriteError(w, http.StatusBadGateway, model.CodeUpstreamError, "通過予測を取得できませんでした")
		return
	}

	WriteJSON(w, http.StatusOK, response)
}

func parsePassesQuery(params url.Values) (service.PassesQuery, error) {
	if err := checkQuery(params, "city", "visibility", "minElevation", "sort", "page", "limit"); err != nil {
		return service.PassesQuery{}, err
	}

	location, ok := findLocation(params.Get("city"))
	if !ok {
		return service.PassesQuery{}, errors.New("cityに対応する都市IDを指定してください")
	}

	query := service.PassesQuery{
		Location: location,
		SortKey:  service.SortByStart,
		Page:     1,
		Limit:    10,
	}

	if params.Has("visibility") {
		for _, value := range strings.Split(params.Get("visibility"), ",") {
			if !slices.Contains(visibilityValues, value) {
				return service.PassesQuery{}, errors.New("visibilityはgood, fair, poor, unknownのいずれかを指定してください")
			}
			query.Visibilities = append(query.Visibilities, value)
		}
	}

	if params.Has("minElevation") {
		value, err := strconv.ParseFloat(params.Get("minElevation"), 64)
		if err != nil || math.IsNaN(value) || value < 0 || value > 90 {
			return service.PassesQuery{}, errors.New("minElevationは0から90の間の数値を指定してください")
		}
		query.MinElevation = &value
	}
	if params.Has("sort") {
		key, descending := strings.CutPrefix(params.Get("sort"), "-")
		if !slices.Contains(sortKeys, key) {
			return service.PassesQuery{}, errors.New("sortはstart・maxElevation・cloudCoverのどれかで指定してください(先頭に-を付けると降順)")
		}
		query.SortKey = key
		query.Descending = descending
	}

	if params.Has("page") {
		value, err := strconv.Atoi(params.Get("page"))
		if err != nil || value < 1 {
			return service.PassesQuery{}, errors.New("pageは1以上の整数で指定してください")
		}
		query.Page = value
	}

	if params.Has("limit") {
		value, err := strconv.Atoi(params.Get("limit"))
		if err != nil || value < 1 || value > 50 {
			return service.PassesQuery{}, errors.New("limitは1〜50の整数で指定してください")
		}
		query.Limit = value
	}

	return query, nil

}

func findLocation(id string) (model.Location, bool) {
	for _, location := range model.Locations {
		if location.ID == id {
			return location, true
		}
	}
	return model.Location{}, false
}
