package interceptor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Interceptor struct {
	logger *zap.Logger
}

func NewInterceptor(logger *zap.Logger) *Interceptor {
	return &Interceptor{
		logger: logger,
	}
}

func (i *Interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	i.logger.Info(fmt.Sprintf("Request Info: %s", info.FullMethod))
	if msg, ok := req.(proto.Message); ok {
		jsonData, err := protojson.Marshal(msg)
		if err == nil {
			i.logger.Info(fmt.Sprintf("Request JSON: %s", jsonData))
		} else {
			i.logger.Error("Failed to marshal request", zap.Error(err))
		}
	}

	res, err := handler(ctx, req)
	if err != nil {
		i.logger.Error("Handler Failed", zap.Error(err))
	} else {
		if msg, ok := res.(proto.Message); ok {
			jsonData, err := protojson.Marshal(msg)
			if err == nil {
				i.logger.Info(fmt.Sprintf("Response JSON: %s", jsonData))
			} else {
				i.logger.Error("Failed to marshal response", zap.Error(err))
			}
		}
	}
	return res, err
}
