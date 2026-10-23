package handler

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/service"
)

func IssNowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		WriteError(
			w,
			http.StatusMethodNotAllowed,
			model.CodeMethodNotAllowed,
			"GETを使用してください",
		)
		return
	}

	// serviceからISSの現在位置を取得
	position, err := service.GetISSPosition(r.Context())
	if err != nil {
		// ブラウザ側のキャンセルは外部APIのタイムアウトと区別
		if errors.Is(err, context.Canceled) {
			return
		}

		var netErr net.Error
		if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()) {

			WriteError(
				w,
				http.StatusGatewayTimeout,
				model.CodeUpstreamTimeout,
				"ISSの位置情報の取得がタイムアウトしました",
			)
			return
		}

		log.Printf("ISS位置情報の取得エラー: %v", err)

		WriteError(
			w,
			http.StatusBadGateway,
			model.CodeUpstreamError,
			"ISSの位置情報を取得できませんでした",
		)
		return
	}

	// 正常時は共通関数でJSONを返す
	WriteJSON(w, http.StatusOK, position)
}
