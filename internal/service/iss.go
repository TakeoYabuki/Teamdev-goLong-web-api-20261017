package service

import (
	"context"
	"sync"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/client"
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const issCacheDuration = 5 * time.Second

var (
	issCacheMu sync.Mutex //複数の処理が同時にキャッシュを書き換えないようにする。
	cachedISS  *model.ISSPosition
	cachedAt   time.Time
)

func GetISSPosition(ctx context.Context) (*model.ISSPosition, error) {
	return getISSPosition(ctx, client.FetchISSPosition)
}

func getISSPosition(
	ctx context.Context,
	fetch func(context.Context) (*model.ISSPosition, error),
) (*model.ISSPosition, error) {

	issCacheMu.Lock()
	defer issCacheMu.Unlock()

	// 5秒以内ならキャッシュを返す
	if cachedISS != nil && time.Since(cachedAt) < issCacheDuration {
		position := *cachedISS
		return &position, nil
	}

	// キャッシュがなければ取得する
	position, err := fetch(ctx)
	if err != nil {
		return nil, err
	}

	// キャッシュを更新
	cachedISS = position
	cachedAt = time.Now()

	result := *position
	return &result, nil
}
