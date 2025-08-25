package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const CtxKeyRequestId = "request_id"

func generateUUID() string {
	return uuid.New().String()
}

func SetRequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CtxKeyRequestId, generateUUID())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestId(r *http.Request) string {
	requestId, ok := r.Context().Value(CtxKeyRequestId).(string)
	if !ok {
		log.Println("Request ID not found in context")
	}
	return requestId
}
