package util

import (
	res "backend/app/response"
	"net/http"
)

func MethodRouter(method map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := method[r.Method]; ok {
			h(w, r)
		} else {
			res.WriteJsonError(w, "許可されていないメソッドです。", http.StatusMethodNotAllowed)
		}
	}
}
