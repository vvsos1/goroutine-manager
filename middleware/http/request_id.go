package http

import (
	"context"
	"net/http"
	"worker-manager/middleware"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.RequestIdKey, middleware.GenerateRequestId(middleware.HttpPrefix))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
