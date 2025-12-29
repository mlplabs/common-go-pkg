package request

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func GetOffsetLimit(r *http.Request) (int, int) {
	varOffset := r.URL.Query().Get("o")
	varLimit := r.URL.Query().Get("l")

	offset, err := strconv.Atoi(varOffset)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(varLimit)
	if err != nil {
		limit = 0
	}
	return offset, limit
}

func GetParamID(r *http.Request) int64 {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0
	}
	return int64(id)
}
