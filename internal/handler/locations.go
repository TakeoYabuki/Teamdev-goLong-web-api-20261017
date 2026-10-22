package handler

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

// 都市検索API専用のレスポンス
type locationsResponse struct {
	Query      string           `json:"query"`
	Results    []model.Location `json:"results"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int              `json:"total"`
	TotalPages int              `json:"totalPages"`
}

func LocationHandler(w http.ResponseWriter, r *http.Request) {
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
	// クエリパラメータを取得する /api/locations/search?q=鹿児島
	params := r.URL.Query()
	query := strings.TrimSpace(params.Get("q")) // 空白カット

	// qの長さ
	if params.Has("q") {
		length := utf8.RuneCountInString(query) // 「鹿児島」の文字数を数える 日本語なのでlenはダメ

		if length < 1 || length > 50 {
			WriteError(
				w,
				http.StatusBadRequest,
				model.CodeInvalidParameter,
				"qは前後の空白を除いて1～50文字で指定してください",
			)
			return
		}
	}

	// ページ番号を取得 page>1
	page := 1
	if params.Has("page") {
		value, err := strconv.Atoi(params.Get("page"))

		if err != nil || value < 1 {
			WriteError(
				w,
				http.StatusBadRequest,
				model.CodeInvalidParameter,
				"pageは1以上の整数で指定してください",
			)
			return
		}
		page = value
	}

	// 1ページの件数を取得 limit 1-50
	limit := 10
	if params.Has("limit") {
		value, err := strconv.Atoi(params.Get("limit"))

		if err != nil || value < 1 || value > 50 {
			WriteError(
				w,
				http.StatusBadRequest,
				model.CodeInvalidParameter,
				"limitは1～50の整数で指定してください",
			)
			return
		}
		limit = value
	}

	// 固定都市を登録順に検索する
	matches := make([]model.Location, 0) // スライスを作る

	for _, location := range model.Locations {
		if query == "" || strings.Contains(location.Name, query) { // 鹿児島市と鹿児島を確認
			matches = append(matches, location)
		}
	}
	// 結果の件数とページ数を計算する
	total := len(matches)
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++ // math.ceil使うと型変換がいるので
	}

	// 指定ページのデータを取り出す 該当なしは空配列
	result := make([]model.Location, 0)

	if page <= totalPages {
		start := (page - 1) * limit
		end := start + limit
		if end > total {
			end = total
		}
		result = matches[start:end]
	}

	// Jsonで返す
	WriteJSON(w, http.StatusOK, locationsResponse{
		Query:      query,
		Results:    result,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}
