package middleware

import (
	"context"
	"log"

	"github.com/google/uuid"
)

const RequestIdKey = "request_id"

type prefix string

const (
	GrpcPrefix prefix = "grpc-"
	HttpPrefix prefix = "http-"
)

func GenerateRequestId(prefix prefix) string {
	return string(prefix) + uuid.New().String()
}

func GetRequestId(ctx context.Context) string {
	if reqId, ok := ctx.Value(RequestIdKey).(string); ok {
		return reqId
	}
	log.Println("request id not found in context")
	return ""
}
