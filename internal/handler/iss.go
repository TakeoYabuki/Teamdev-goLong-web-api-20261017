package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
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

	// Where the ISS at? のモックデータ
	jst := time.FixedZone("JST", 9*60*60)

	position := model.ISSPosition{
		Lat:         -14.97657049433,
		Lon:         -167.91305342777,
		AltitudeKm:  420.94118518006,
		VelocityKmh: 27578.208454235,
		PositionAt:  time.Date(2026, 10, 19, 15, 30, 0, 0, jst),
		FetchedAt:   time.Date(2026, 10, 19, 15, 30, 5, 0, jst),
	}

	w.Header().Set("Content-Type", "application/json")

	//JSON形成
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(position); err != nil {
		log.Printf("JSONエンコードエラー: %v", err)
	}
}
