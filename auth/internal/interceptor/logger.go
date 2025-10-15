package interceptor

import (
	"auth/internal/logger"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("request", req))
	}

	logger.Info("successful", zap.String("method", info.FullMethod), zap.Any("response", resp), zap.Duration("duration", time.Since(now)))

	return resp, err
}
