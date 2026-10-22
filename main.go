package main

import (
	"log"
	"net/http"

	//hanlderフォルダ内の関数呼び出し用。
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	//URLと処理を対応付ける
	mux.HandleFunc("/api/health", handler.HealthHandler)
	//internal/handler/iss.goへのルーティン
	mux.HandleFunc("/api/iss/now", handler.IssNowHandler)

	log.Println("サーバを起動します: http://localhost:8000")

	if err := http.ListenAndServe(":8000", handler.WithCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
