package main

import (
	"log"
	"net/http"
	"os"

	//hanlderフォルダ内の関数呼び出し用。
	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	//URLと処理を対応付ける
	mux.HandleFunc("/api/health", handler.HealthHandler)
	//internal/handler/iss.goへのルーティン
	mux.HandleFunc("/api/iss/now", handler.IssNowHandler)
	//internal/handler/locations.goへのルーティン
	mux.HandleFunc("/api/locations/search", handler.LocationHandler)
	mux.HandleFunc("/api/iss/passes", handler.PassesHandler)
	mux.HandleFunc("/api/iss/passes/{id}", handler.PassDetailHandler)
	mux.HandleFunc("/", handler.NotFoundHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("サーバを起動します: http://localhost:" + port)

	if err := http.ListenAndServe("127.0.0.1:"+port, handler.WithCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
