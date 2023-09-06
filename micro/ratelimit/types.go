package ratelimit

import (
	"context"
	"google.golang.org/grpc"
)

type rejectStrategy func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)

var defaultRejectStrategy rejectStrategy = func(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// 设置一个标志位
	ctx = context.WithValue(ctx, "limited", true)
	return handler(ctx, req)
}
