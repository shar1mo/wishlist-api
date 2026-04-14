package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseInt64Param(r *http.Request, key string) (int64, error) {
	value := chi.URLParam(r, key)
	return strconv.ParseInt(value, 10, 64)
}