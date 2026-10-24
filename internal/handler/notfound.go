package handler

import (
	"net/http"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusNotFound, model.CodeNotFound, "指定したAPIは存在しません")
}
