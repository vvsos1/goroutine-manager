package grpc

import (
	"context"
	"worker-manager/middleware"

	"google.golang.org/grpc"
)

func RequestId(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	newCtx := context.WithValue(ctx, middleware.RequestIdKey, middleware.GenerateRequestId(middleware.GrpcPrefix))
	return handler(newCtx, req)
}
