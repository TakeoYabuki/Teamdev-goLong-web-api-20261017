package service

import (
	"slices"
	"testing"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

var testNow = time.Date(2026, 10, 27, 18, 0, 0, 0, jst)

func newTestPass(id string, startAfter time.Duration, maxElevation float64, magnitude *float64) model.Pass {
	start := testNow.Add(startAfter)
	return model.Pass{
		ID:        id,
		Start:     model.PassPoint{Time: start, Direction: "北西", Elevation: 10},
		Max:       model.PassPoint{Time: start.Add(3 * time.Minute), Direction: "北", Elevation: maxElevation},
		End:       model.PassPoint{Time: start.Add(6 * time.Minute), Direction: "北東", Elevation: 10},
		Magnitude: magnitude,
	}
}

func testPasses() []model.Pass {
	return []model.Pass{
		newTestPass("before", -time.Hour, 62, ptr(-3)),
		newTestPass("good", time.Hour, 62, ptr(-3)),
		newTestPass("fair", 25*time.Hour, 35, ptr(-1.5)),
		newTestPass("poor", 49*time.Hour, 15, ptr(-3)),
		newTestPass("unknown", 60*time.Hour, 50, ptr(-2.5)),
		newTestPass("after", 72*time.Hour, 62, ptr(-3)),
	}
}

func testForecast() *model.CloudForecast {
	return &model.CloudForecast{
		Hours: []model.HourlyCloud{
			{Time: testNow.Add(time.Hour), CloudCover: ptr(20)},
			{Time: testNow.Add(25 * time.Hour), CloudCover: ptr(45)},
			{Time: testNow.Add(49 * time.Hour), CloudCover: ptr(10)},
		},
		FetchedAt: testNow,
	}
}

func passIDs(passes []model.Pass) []string {
	ids := []string{}
	for _, pass := range passes {
		ids = append(ids, pass.ID)
	}
	return ids
}

func TestBuildPassesResponseWindowAndJudge(t *testing.T) {
	query := PassesQuery{SortKey: SortByStart, Page: 1, Limit: 10}

	got := buildPassesResponse(query, testPasses(), testNow, testForecast(), testNow)

	if want := []string{"good", "fair", "poor", "unknown"}; !slices.Equal(passIDs(got.Passes), want) {
		t.Fatalf("72時間の通過 = %v, want %v", passIDs(got.Passes), want)
	}
	for i, want := range []string{"good", "fair", "poor", "unknown"} {
		if got.Passes[i].Visibility != want {
			t.Errorf("%s の判定 = %s, want %s", got.Passes[i].ID, got.Passes[i].Visibility, want)
		}
	}
	if got.Passes[0].CloudCover == nil || *got.Passes[0].CloudCover != 20 {
		t.Errorf("雲量が付いていません: %v", got.Passes[0].CloudCover)
	}
	if got.Passes[0].Message == "" || len(got.Passes[0].Reasons) == 0 {
		t.Error("案内文か理由が空です")
	}
	if !got.From.Equal(testNow) || !got.To.Equal(testNow.Add(72*time.Hour)) {
		t.Errorf("期間 = %v 〜 %v", got.From, got.To)
	}
	if want := (model.PassSummary{Total: 4, Good: 1, Fair: 1, Poor: 1, Unknown: 1}); got.Summary != want {
		t.Errorf("summary = %+v, want %+v", got.Summary, want)
	}
	if got.NextCandidate == nil || got.NextCandidate.ID != "good" {
		t.Errorf("nextCandidate = %v, want good", got.NextCandidate)
	}
}

func TestBuildPassesResponseFilterKeepsSummaryAndNextCandidate(t *testing.T) {
	minElevation := 40.0
	query := PassesQuery{Visibilities: []string{"good", "unknown"}, MinElevation: &minElevation, SortKey: SortByStart, Page: 1, Limit: 10}

	got := buildPassesResponse(query, testPasses(), testNow, testForecast(), testNow)

	if want := []string{"good", "unknown"}; !slices.Equal(passIDs(got.Passes), want) {
		t.Errorf("絞り込み後 = %v, want %v", passIDs(got.Passes), want)
	}
	if got.Total != 2 || got.Summary.Total != 4 {
		t.Errorf("total = %d(want 2)、summary.total = %d(want 4)", got.Total, got.Summary.Total)
	}

	onlyPoor := PassesQuery{Visibilities: []string{"poor"}, SortKey: SortByStart, Page: 1, Limit: 10}
	got = buildPassesResponse(onlyPoor, testPasses(), testNow, testForecast(), testNow)
	if got.NextCandidate == nil || got.NextCandidate.ID != "good" {
		t.Errorf("絞り込みに関係なく最も早い通過のはずが、nextCandidate = %v", got.NextCandidate)
	}
}

func TestBuildPassesResponseSort(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		descending bool
		want       []string
	}{
		{"開始時刻の降順", SortByStart, true, []string{"unknown", "poor", "fair", "good"}},
		{"最大仰角の降順", SortByMaxElevation, true, []string{"good", "unknown", "fair", "poor"}},
		{"雲量の昇順でnullは最後", SortByCloudCover, false, []string{"poor", "good", "fair", "unknown"}},
		{"雲量の降順でもnullは最後", SortByCloudCover, true, []string{"fair", "good", "poor", "unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := PassesQuery{SortKey: tt.key, Descending: tt.descending, Page: 1, Limit: 10}

			got := buildPassesResponse(query, testPasses(), testNow, testForecast(), testNow)

			if !slices.Equal(passIDs(got.Passes), tt.want) {
				t.Errorf("並び = %v, want %v", passIDs(got.Passes), tt.want)
			}
		})
	}
}

func TestBuildPassesResponseUnknownCloudIsLast(t *testing.T) {
	passes := []model.Pass{
		newTestPass("noCloud", time.Hour, 50, ptr(-3)),
		newTestPass("cloudy", 2*time.Hour, 50, ptr(-3)),
		newTestPass("clear", 3*time.Hour, 50, ptr(-3)),
	}
	forecast := &model.CloudForecast{
		Hours: []model.HourlyCloud{
			{Time: testNow.Add(2 * time.Hour), CloudCover: ptr(80)},
			{Time: testNow.Add(3 * time.Hour), CloudCover: ptr(10)},
		},
		FetchedAt: testNow,
	}

	tests := []struct {
		descending bool
		want       []string
	}{
		{false, []string{"clear", "cloudy", "noCloud"}},
		{true, []string{"cloudy", "clear", "noCloud"}},
	}
	for _, tt := range tests {
		query := PassesQuery{SortKey: SortByCloudCover, Descending: tt.descending, Page: 1, Limit: 10}

		got := buildPassesResponse(query, passes, testNow, forecast, testNow)

		if !slices.Equal(passIDs(got.Passes), tt.want) {
			t.Errorf("降順=%v の並び = %v, want %v", tt.descending, passIDs(got.Passes), tt.want)
		}
	}
}

func TestBuildPassesResponseTieBreak(t *testing.T) {
	passes := []model.Pass{
		newTestPass("b", 3*time.Hour, 50, ptr(-3)),
		newTestPass("c", 2*time.Hour, 50, ptr(-3)),
		newTestPass("a", 3*time.Hour, 50, ptr(-3)),
	}
	query := PassesQuery{SortKey: SortByMaxElevation, Descending: true, Page: 1, Limit: 10}

	got := buildPassesResponse(query, passes, testNow, nil, testNow)

	if want := []string{"c", "a", "b"}; !slices.Equal(passIDs(got.Passes), want) {
		t.Errorf("同じ値の並び = %v, want %v(開始時刻、IDの順)", passIDs(got.Passes), want)
	}
}

func TestBuildPassesResponsePaging(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		limit          int
		want           []string
		wantTotalPages int
	}{
		{"2ページ目", 2, 3, []string{"unknown"}, 2},
		{"範囲外のページは空", 3, 3, []string{}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := PassesQuery{SortKey: SortByStart, Page: tt.page, Limit: tt.limit}

			got := buildPassesResponse(query, testPasses(), testNow, testForecast(), testNow)

			if !slices.Equal(passIDs(got.Passes), tt.want) {
				t.Errorf("通過 = %v, want %v", passIDs(got.Passes), tt.want)
			}
			if got.TotalPages != tt.wantTotalPages {
				t.Errorf("totalPages = %d, want %d", got.TotalPages, tt.wantTotalPages)
			}
		})
	}
}

func TestBuildPassesResponseNoPassesAndNoWeather(t *testing.T) {
	query := PassesQuery{SortKey: SortByStart, Page: 1, Limit: 10}

	got := buildPassesResponse(query, nil, testNow, nil, testNow)

	if got.Passes == nil || len(got.Passes) != 0 {
		t.Errorf("passes は空の配列のはずが %v", got.Passes)
	}
	if got.NextCandidate != nil || got.TotalPages != 0 || got.WeatherFetchedAt != nil {
		t.Errorf("nextCandidate=%v totalPages=%d weatherFetchedAt=%v", got.NextCandidate, got.TotalPages, got.WeatherFetchedAt)
	}

	got = buildPassesResponse(query, testPasses(), testNow, nil, testNow)
	for _, pass := range got.Passes {
		if pass.Visibility != model.VisibilityPoor && pass.Visibility != model.VisibilityUnknown {
			t.Errorf("天気がないのに %s が %s になりました", pass.ID, pass.Visibility)
		}
	}
}
