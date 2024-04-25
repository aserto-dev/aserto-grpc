package tracing

import (
	"context"
	"time"

	"github.com/aserto-dev/aserto-grpc/grpcutil"
	"github.com/aserto-dev/header"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type TracingMiddleware struct {
}

func NewTracingMiddleware() *TracingMiddleware {
	return &TracingMiddleware{}
}

var _ grpcutil.Middleware = &TracingMiddleware{}

func (m *TracingMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		apiLogger := zerolog.Ctx(ctx).With().
			Fields(header.KnownContextValueStrings(ctx)).
			Logger()

		apiLogger.Trace().Interface("request", req).Msg("grpc call start")

		newCtx := apiLogger.WithContext(ctx)

		start := time.Now()
		result, err := handler(newCtx, req)
		apiLogger.Trace().Dur("duration-ms", time.Since(start)).Msg("grpc call end")

		return result, err
	}
}

func (m *TracingMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()

		apiLogger := zerolog.Ctx(ctx)

		apiLogger.Trace().
			Fields(header.KnownContextValueStrings(ctx)).
			Msg("grpc stream call")

		newCtx := apiLogger.WithContext(ctx)

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
