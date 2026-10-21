package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("JSONの書き出しに失敗しました: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, code string, message string) {
	WriteJSON(w, status, model.ErrorResponse{
		Error: model.ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}
