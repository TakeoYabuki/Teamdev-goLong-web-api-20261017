package service

import (
	"context"
	"testing"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestSnapshotStore(t *testing.T) {
	savedAt := time.Date(2026, 10, 28, 18, 0, 0, 0, jst)
	store := newSnapshotStore()
	store.save("tokyo", []model.Pass{{ID: "tokyo_1", Message: "最初の一覧"}}, savedAt)

	if pass, found := store.get("tokyo", "tokyo_1", savedAt.Add(time.Hour)); !found || pass.Message != "最初の一覧" {
		t.Errorf("保存した通過が取れません: found=%v pass=%+v", found, pass)
	}
	if _, found := store.get("osaka", "tokyo_1", savedAt); found {
		t.Error("別の都市では見つからないはずです")
	}
	if _, found := store.get("tokyo", "tokyo_2", savedAt); found {
		t.Error("保存していないIDは見つからないはずです")
	}

	store.save("tokyo", []model.Pass{{ID: "tokyo_1", Message: "あとの一覧"}}, savedAt.Add(time.Hour))
	if pass, _ := store.get("tokyo", "tokyo_1", savedAt.Add(time.Hour)); pass.Message != "最初の一覧" {
		t.Errorf("同じIDは上書きしないはずが、%q になりました", pass.Message)
	}

	if _, found := store.get("tokyo", "tokyo_1", savedAt.Add(24*time.Hour)); found {
		t.Error("24時間たったら見つからないはずです")
	}
}

func TestSnapshotStoreRemovesExpired(t *testing.T) {
	savedAt := time.Date(2026, 10, 28, 18, 0, 0, 0, jst)
	store := newSnapshotStore()
	store.save("tokyo", []model.Pass{{ID: "tokyo_1"}}, savedAt)

	store.save("tokyo", []model.Pass{{ID: "tokyo_2"}}, savedAt.Add(25*time.Hour))

	if len(store.snapshots) != 1 {
		t.Errorf("期限切れが消えていません: 残り %d 件", len(store.snapshots))
	}
}

func TestGetPassesSavesSnapshotsForDetail(t *testing.T) {
	now := time.Now()
	newPass := func(id string, startAfter time.Duration) model.Pass {
		start := now.Add(startAfter)
		return model.Pass{
			ID:    id,
			Start: model.PassPoint{Time: start, Direction: "北西"},
			Max:   model.PassPoint{Time: start.Add(3 * time.Minute), Elevation: 62},
		}
	}
	passes := []model.Pass{
		newPass("kagoshima_1", time.Hour),
		newPass("kagoshima_2", 2*time.Hour),
		newPass("kagoshima_3", 3*time.Hour),
	}
	passesCache.entries["kagoshima"] = &cacheEntry[[]model.Pass]{value: passes, storedAt: now, stored: true}
	forecastCache.entries["kagoshima"] = &cacheEntry[*model.CloudForecast]{value: &model.CloudForecast{FetchedAt: now}, storedAt: now, stored: true}
	t.Cleanup(func() {
		delete(passesCache.entries, "kagoshima")
		delete(forecastCache.entries, "kagoshima")
		passSnapshots = newSnapshotStore()
	})

	query := PassesQuery{Location: model.Location{ID: "kagoshima"}, SortKey: SortByStart, Page: 2, Limit: 1}
	if _, err := GetPasses(context.Background(), query); err != nil {
		t.Fatalf("一覧の取得に失敗しました: %v", err)
	}

	tests := []struct {
		id        string
		wantFound bool
	}{
		{"kagoshima_1", true},
		{"kagoshima_2", true},
		{"kagoshima_3", false},
	}
	for _, tt := range tests {
		detail, found := GetPassDetail("kagoshima", tt.id)
		if found != tt.wantFound {
			t.Errorf("%s が見つかったか = %v, want %v", tt.id, found, tt.wantFound)
		}
		if found && (detail.Visibility == "" || detail.Message == "") {
			t.Errorf("%s の詳細に判定と案内文が入っていません", tt.id)
		}
	}
}
