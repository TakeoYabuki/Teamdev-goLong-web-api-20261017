package handler

import (
	"encoding/json"
	"net/http"
)

// サーバーの動作確認
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		// GETじゃなかったら、GETだけだよと伝えて405へ
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "GETを使用してください",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "OK",
	})
}
