package service

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newTestCache(clock *time.Time) *cityCache[string] {
	cache := newCityCache[string]()
	cache.now = func() time.Time { return *clock }
	return cache
}

func TestCityCacheReusesWithinAnHour(t *testing.T) {
	clock := time.Date(2026, 10, 27, 18, 0, 0, 0, jst)
	cache := newTestCache(&clock)
	fetchCount := 0
	fetch := func() (string, error) {
		fetchCount++
		return "tokyoの通過", nil
	}

	first, firstStoredAt, _ := cache.get("tokyo", fetch)
	clock = clock.Add(59 * time.Minute)
	second, secondStoredAt, _ := cache.get("tokyo", fetch)

	if fetchCount != 1 {
		t.Errorf("取得の回数 = %d, want 1", fetchCount)
	}
	if first != second || !firstStoredAt.Equal(secondStoredAt) {
		t.Errorf("1回目と2回目が違います: %q %v / %q %v", first, firstStoredAt, second, secondStoredAt)
	}
}

func TestCityCacheFetchesAgainAfterAnHour(t *testing.T) {
	clock := time.Date(2026, 10, 27, 18, 0, 0, 0, jst)
	cache := newTestCache(&clock)
	fetchCount := 0
	fetch := func() (string, error) {
		fetchCount++
		return "tokyoの通過", nil
	}

	cache.get("tokyo", fetch)
	clock = clock.Add(time.Hour)
	_, storedAt, _ := cache.get("tokyo", fetch)

	if fetchCount != 2 {
		t.Errorf("取得の回数 = %d, want 2", fetchCount)
	}
	if !storedAt.Equal(clock) {
		t.Errorf("保存した時刻 = %v, want %v", storedAt, clock)
	}
}

func TestCityCacheDoesNotStoreFailure(t *testing.T) {
	clock := time.Date(2026, 10, 27, 18, 0, 0, 0, jst)
	cache := newTestCache(&clock)
	fetchCount := 0
	failing := func() (string, error) {
		fetchCount++
		return "", errors.New("N2YOの失敗")
	}

	if _, _, err := cache.get("tokyo", failing); err == nil {
		t.Fatal("失敗がそのまま返るはずが、エラーになりませんでした")
	}
	cache.get("tokyo", failing)

	if fetchCount != 2 {
		t.Errorf("失敗を保存せず取り直すはずが、取得の回数 = %d", fetchCount)
	}
}

func TestCityCacheSeparatesCities(t *testing.T) {
	clock := time.Date(2026, 10, 27, 18, 0, 0, 0, jst)
	cache := newTestCache(&clock)

	tokyo, _, _ := cache.get("tokyo", func() (string, error) { return "tokyoの通過", nil })
	osaka, _, _ := cache.get("osaka", func() (string, error) { return "osakaの通過", nil })

	if tokyo != "tokyoの通過" || osaka != "osakaの通過" {
		t.Errorf("都市ごとに分かれていません: %q %q", tokyo, osaka)
	}
}

func TestCityCacheSharesConcurrentFetch(t *testing.T) {
	cache := newCityCache[string]()
	var fetchCount atomic.Int32
	fetch := func() (string, error) {
		fetchCount.Add(1)
		time.Sleep(50 * time.Millisecond)
		return "tokyoの通過", nil
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.get("tokyo", fetch)
		}()
	}
	wg.Wait()

	if got := fetchCount.Load(); got != 1 {
		t.Errorf("同時に10回呼んだときの取得の回数 = %d, want 1", got)
	}
}
