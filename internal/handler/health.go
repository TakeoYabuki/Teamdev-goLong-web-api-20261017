package handler

import (
	"net/http"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		WriteError(w, http.StatusMethodNotAllowed, model.CodeMethodNotAllowed, "GETを使用してください")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}
