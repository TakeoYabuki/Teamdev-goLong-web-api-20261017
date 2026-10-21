package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// モックデータ作成用のstruct
type ISSPosition struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Altitude   float64 `json:"altitude"`
	Velocity   float64 `json:"velocity"`
	PositionAt int64   `json:"positionAt"`
}

func IssNowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "GETを使用してください", http.StatusMethodNotAllowed)
		return
	}

	// Where the ISS at? のモックデータ
	position := ISSPosition{
		Latitude:   -14.97657049433,
		Longitude:  -167.91305342777,
		Altitude:   420.94118518006,
		Velocity:   27578.208454235,
		PositionAt: 1792494475,
	}

	w.Header().Set("Content-Type", "application/json")

	//JSON形成
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(position); err != nil {
		log.Printf("JSONエンコードエラー: %v", err)
	}
}
