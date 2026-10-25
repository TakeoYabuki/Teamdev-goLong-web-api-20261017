package service

import (
	"context"
	"testing"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestGetISSPosition_Cache(t *testing.T) {
	// 他のテストのキャッシュが残らないよう初期化
	issCacheMu.Lock()
	cachedISS = nil
	cachedAt = time.Time{}
	issCacheMu.Unlock()

	// 外部APIの代わりになる関数
	fetchCount := 0

	fakeFetch := func(ctx context.Context) (*model.ISSPosition, error) {
		fetchCount++

		return &model.ISSPosition{
			Lat: 35.68,
			Lon: 139.69,
		}, nil
	}

	// 1回目：取得処理が実行される
	first, err := getISSPosition(context.Background(), fakeFetch)
	if err != nil {
		t.Fatalf("1回目の取得に失敗: %v", err)
	}

	// 2回目：キャッシュが使われる
	second, err := getISSPosition(context.Background(), fakeFetch)
	if err != nil {
		t.Fatalf("2回目の取得に失敗: %v", err)
	}

	// 取得処理が1回しか実行されていないことを確認
	if fetchCount != 1 {
		t.Errorf("取得回数は1回を期待しましたが、実際は%d回", fetchCount)
	}

	// 1回目と2回目のデータが同じことを確認
	if first.Lat != second.Lat || first.Lon != second.Lon {
		t.Errorf("キャッシュの値が一致しません")
	}
}
