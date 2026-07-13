package gapi

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startedAt)

		st := status.Convert(err)
		if err != nil {
			log.Printf("[grpc] method=%s code=%s duration=%s error=%v", info.FullMethod, st.Code(), duration, err)
			return resp, err
		}

		log.Printf("[grpc] method=%s code=%s duration=%s", info.FullMethod, st.Code(), duration)
		return resp, nil
	}
}
