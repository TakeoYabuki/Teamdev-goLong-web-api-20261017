package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// サーバーの動作確認
func healthHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	mux := http.NewServeMux()

	//URLと処理を対応付ける
	mux.HandleFunc("/api/health", healthHandler)

	log.Println("サーバを起動します: http://localhost:8000")

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
