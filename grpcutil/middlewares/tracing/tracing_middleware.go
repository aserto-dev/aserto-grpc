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
	logger *zerolog.Logger
}

type tracingHook struct{}

func (h tracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()

	e.Fields(header.KnownContextValueStrings(ctx))

	serviceMethod, ok := grpc.Method(ctx)
	if ok {
		e.Str("method", serviceMethod)

	}

}

func NewTracingMiddleware(logger *zerolog.Logger) *TracingMiddleware {
	return &TracingMiddleware{
		logger: logger,
	}
}

var _ grpcutil.Middleware = &TracingMiddleware{}

func (m *TracingMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := m.logger.Hook(tracingHook{}).With().Interface("request", req).Ctx(ctx).Logger()
		loggerCtx := logger.WithContext(ctx)

		logger.Trace().Msg("grpc call start")

		start := time.Now()
		result, err := handler(loggerCtx, req)
		logger.Trace().Dur("duration-ms", time.Since(start)).Msg("grpc call end")

		return result, err
	}
}

func (m *TracingMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		logger := m.logger.Hook(tracingHook{}).With().Ctx(ctx).Logger()
		loggerCtx := logger.WithContext(ctx)

		logger.Trace().Msg("grpc stream call")

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = loggerCtx

		return handler(srv, wrapped)
	}
}
