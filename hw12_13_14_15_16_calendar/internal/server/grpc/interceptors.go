package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func NewUnaryChainOption(log Logger) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		LoggingUnaryInterceptor(log))
}

func LoggingUnaryInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		peerAddr := ""
		if p, ok := peer.FromContext(ctx); ok {
			peerAddr = p.Addr.String()
		}
		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}
		log.Info(fmt.Sprintf("%s %s %s %d",
			peerAddr, info.FullMethod, code, time.Since(start).Milliseconds()))
		return resp, err
	}
}
