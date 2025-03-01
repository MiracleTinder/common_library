package interceptors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(ctx context.Context,
	req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	requestID := uuid.New().String()
	requestIDField := zap.Fields(zap.String("request_id", requestID))
	logger, err := zap.NewProduction(requestIDField)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed initialize logger")
	}

	logger.Info("request", zap.Time("request_time", time.Now()), zap.String("method", info.FullMethod))
	resStart := time.Now()
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error("Failed create response")
		return nil, status.Error(codes.Internal, "Failed create response")
	}

	logger.Info("response", zap.Duration("response_time", time.Since(resStart)),
		zap.String("method", info.FullMethod))

	return resp, err
}
